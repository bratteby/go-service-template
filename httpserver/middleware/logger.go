package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap/zapcore"

	"github.com/bratteby/go-service-template/logging"
)

// RequestLoggerOptions contains the middleware configuration.
type RequestLoggerOptions struct {
	// Enable verbose logging
	Verbose bool
}

var defaultLoggerOptions = RequestLoggerOptions{
	Verbose: false,
}

// RequestLogger is an http middleware to log http requests and responses.
func RequestLogger(logger logging.Logger, opts *RequestLoggerOptions) func(next http.Handler) http.Handler {
	if opts != nil {
		defaultLoggerOptions = *opts
	}

	var f chimiddleware.LogFormatter = &requestLogger{logger}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				entry = f.NewLogEntry(r)
				ww    = chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
				buf   = newLimitBuffer(512)
				t1    = time.Now()
			)

			ww.Tee(buf)

			defer func() {
				var respBody []byte
				if ww.Status() >= 400 {
					respBody, _ = ioutil.ReadAll(buf)
				}
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), respBody)
			}()

			next.ServeHTTP(ww, chimiddleware.WithLogEntry(r, entry))
		}

		return http.HandlerFunc(fn)
	}
}

// requestLogger implements the middleware.LogFormatter interface.
type requestLogger struct {
	Logger logging.Logger
}

// NewLogEntry creates a new LogEntry for the request.
func (l *requestLogger) NewLogEntry(r *http.Request) chimiddleware.LogEntry {
	entry := &RequestLoggerEntry{}
	msg := fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path)
	entry.Logger = l.Logger.With(requestLogFields(r)...)

	if defaultLoggerOptions.Verbose {
		entry.Logger.Info(msg)
	}

	return entry
}

type RequestLoggerEntry struct {
	Logger logging.Logger
	msg    string
}

func (l *RequestLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	msg := fmt.Sprintf("Response: %d %s", status, statusLabel(status))
	if l.msg != "" {
		msg = fmt.Sprintf("%s - %s", msg, l.msg)
	}

	responseLog := respLog{
		Status:       status,
		BytesWritten: bytes,
		Elapsed:      elapsed.Milliseconds(),
	}

	// Include response header, as well for error status codes (>400) we include
	// the response body so we may inspect the log message sent back to the client.
	if status >= 400 {
		body, _ := extra.([]byte)
		responseLog.Body = string(body)
	}

	if len(header) > 0 {
		responseLog.Header = getHeaderLogField(header)
	}

	logLevel := statusLevel(status)
	responseLogField := logging.Object("httpResponse", &responseLog)

	switch logLevel {
	case logging.InfoLevel:
		l.Logger.InfoWith(msg, responseLogField)
	case logging.ErrorLevel:
		l.Logger.ErrorWith(msg, responseLogField)
	}

}

func (l *RequestLoggerEntry) Panic(v interface{}, stack []byte) {
	stacktrace := string(stack)

	l.Logger = l.Logger.With(
		logging.String("stacktrace", stacktrace),
		logging.String("panic", fmt.Sprintf("%+v", v)),
	)

	l.msg = fmt.Sprintf("%+v", v)
}

// headerFields contains values from headers in a request or response.
type headerFields map[string]string

// reqLog defines values to be logged from a http request.
type reqLog struct {
	RequestURL    string
	RequestMethod string
	RequestPath   string
	RemoteIP      string
	Proto         string
	ReqID         string

	// Maybe
	Header headerFields
}

// reqLog defines values to be logged from a http response.
type respLog struct {
	Status       int
	BytesWritten int
	Elapsed      int64

	// Maybe
	Body   string
	Header headerFields
}

func (f *reqLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("requestURL", f.RequestURL)
	enc.AddString("requestMethod", f.RequestMethod)
	enc.AddString("requestPath", f.RequestPath)
	enc.AddString("remoteIP", f.RemoteIP)
	enc.AddString("proto", f.Proto)

	if f.ReqID != "" {
		enc.AddString("requestID", f.ReqID)
	}

	if f.Header != nil {
		enc.AddObject("header", &f.Header)
	}

	return nil
}

func (f *respLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status", f.Status)
	enc.AddInt("bytesWritten", f.BytesWritten)
	enc.AddInt64("elapsed", f.Elapsed)

	if f.Body != "" {
		enc.AddString("body", f.Body)
	}

	if f.Header != nil {
		enc.AddObject("header", &f.Header)
	}

	return nil
}

func (hm *headerFields) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range *hm {
		enc.AddString(k, v)
	}

	return nil
}

func requestLogFields(r *http.Request) []logging.Field {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	requestURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	requestFields := []logging.Field{
		logging.Object("httpRequest", &reqLog{
			RequestURL:    requestURL,
			RequestMethod: r.Method,
			RequestPath:   r.URL.Path,
			RemoteIP:      r.RemoteAddr,
			Proto:         r.Proto,
			ReqID:         chimiddleware.GetReqID(r.Context()),
			Header:        getHeaderLogField(r.Header),
		}),
	}

	return requestFields
}

func getHeaderLogField(header http.Header) map[string]string {
	headerField := map[string]string{}
	for k, v := range header {
		k = strings.ToLower(k)
		switch {
		case len(v) == 0:
			continue
		case len(v) == 1:
			headerField[k] = v[0]
		default:
			headerField[k] = fmt.Sprintf("[%s]", strings.Join(v, "], ["))
		}
		if k == "authorization" || k == "cookie" || k == "set-cookie" {
			headerField[k] = "***"
		}

	}
	return headerField
}

func statusLevel(status int) logging.Level {
	switch {
	case status >= 200 && status < 400:
		return logging.InfoLevel
	case status >= 400:
		return logging.ErrorLevel
	default:
		return logging.InfoLevel
	}
}

func statusLabel(status int) string {
	switch {
	case status >= 100 && status < 300:
		return "OK"
	case status >= 300 && status < 400:
		return "Redirect"
	case status >= 400 && status < 500:
		return "Client Error"
	case status >= 500:
		return "Server Error"
	default:
		return "Unknown"
	}
}

///////

// limitBuffer is used to pipe response body information from the
// response writer to a certain limit amount. The idea is to read
// a portion of the response body such as an error response so we
// may log it.
type limitBuffer struct {
	*bytes.Buffer
	limit int
}

func newLimitBuffer(size int) io.ReadWriter {
	return limitBuffer{
		Buffer: bytes.NewBuffer(make([]byte, 0, size)),
		limit:  size,
	}
}

func (b limitBuffer) Write(p []byte) (n int, err error) {
	if b.Buffer.Len() >= b.limit {
		return len(p), nil
	}
	limit := b.limit
	if len(p) < limit {
		limit = len(p)
	}
	return b.Buffer.Write(p[:limit])
}

func (b limitBuffer) Read(p []byte) (n int, err error) {
	return b.Buffer.Read(p)
}
