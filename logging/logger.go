// github.com/bigwhite/experiments/tree/master/uber-zap-advanced-usage/demo1/pkg/log/log.go
package logging

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Also from: https://www.sobyte.net/post/2022-03/uber-zap-advanced-usage/#1-wrapping-zap-to-make-it-more-usable

type Level = zapcore.Level

const (
	InfoLevel  Level = zap.InfoLevel  // 0, default level
	ErrorLevel Level = zap.ErrorLevel // 2
	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

// function variables for all field types in github.com/uber-go/zap/field.go
var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString

	//... ...

	Float64   = zap.Float64
	Float64p  = zap.Float64p
	Float32   = zap.Float32
	Float32p  = zap.Float32p
	Durationp = zap.Durationp

	//... ...

	Int    = zap.Int
	String = zap.String
)

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Infof(template string, args ...any) {
	l.l.Sugar().Infof(template, args...)
}

// InfoWith logs a message with some additional context.
// The variadic keysAndValues are processed in pairs, the first element of the pair is used as the field key and the second as the field value.
func (l *Logger) InfoWith(msg string, keysAndValues ...interface{}) {
	l.l.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Error(err error, fields ...Field) {
	l.l.Error(err.Error(), fields...)
}

// ErrorWith logs a message with some additional context.
// The variadic keysAndValues are processed in pairs, the first element of the pair is used as the field key and the second as the field value.
func (l *Logger) ErrorWith(msg string, keysAndValues ...interface{}) {
	l.l.Sugar().Errorw(msg, keysAndValues...)
}

// Sync calls the underlying Zap's Sync method, flushing any buffered log entries. Applications should take care to call Sync before exiting.
func (l *Logger) Sync() error {
	return l.l.Sync()
}

type Option = zap.Option

var (
	// AddStacktrace configures the Logger to record a stack trace for all messages at or above a given level.
	AddStacktrace = zap.AddStacktrace
)

// Config provides logging configuration.
type Config struct {
	Level         Level // Min logged level.
	WithTimeStamp bool  // Log messages with timestamp.
	Options       []Option
}

// New create a new logger given a writer, min logged level and options.
// The writer is typically nil, in which case os.Stderr is used.
func New(writer io.Writer, conf Config) *Logger {
	if writer == nil {
		writer = os.Stderr
	}

	zapConfig := zap.NewProductionConfig()

	// Format timestamp.
	if conf.WithTimeStamp {
		zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	} else {
		zapConfig.EncoderConfig.TimeKey = ""
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(conf.Level),
	)

	logger := &Logger{
		l:     zap.New(core, conf.Options...),
		level: conf.Level,
	}

	return logger
}
