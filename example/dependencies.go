package example

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out mock_example_repository_test.go . exampleRepository
type exampleRepository interface {
	FindOneByID(ctx context.Context, id uuid.UUID) (Example, error)
	Save(ctx context.Context, ex Example) error
}
