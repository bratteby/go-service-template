package httpserver

import (
	"context"

	"github.com/bratteby/go-service-template/example"
	"github.com/google/uuid"
)

type exampleService interface {
	CreateExample(context.Context, example.ExampleDTO) (example.Example, error)
	GetExampleByID(context.Context, uuid.UUID) (example.Example, error)
}
