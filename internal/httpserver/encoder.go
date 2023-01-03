package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bratteby/go-service-template/internal/example"
	"github.com/bratteby/go-service-template/internal/logging"
)

type encoder struct {
	Logger *logging.Logger
}

// errorResponse will encapsulate errors to be transferred over HTTP.
type errorResponse struct {
	Message string `json:"message"`
}

func (e encoder) respond(
	ctx context.Context,
	w http.ResponseWriter,
	response any,
	statusCode int,
) {
	if response == nil {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		e.Logger.Error(fmt.Errorf("error encoding/writing response %w", err))
	}
}

func (e encoder) error(ctx context.Context, w http.ResponseWriter, err error) {
	e.Logger.ErrorWith(err.Error())

	var (
		apiErr     example.APIError
		statusCode int
		errorMsg   string
	)

	if errors.As(err, &apiErr) {
		statusCode, errorMsg = apiErr.APIError()

	} else {
		statusCode = http.StatusInternalServerError
		errorMsg = "internal error"
	}

	resp := errorResponse{
		Message: errorMsg,
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		e.Logger.Error(fmt.Errorf("error encoding/writing response %w", err))
	}
}
