package example

import (
	"context"
	"fmt"

	"github.com/bratteby/go-service-template/internal/logging"
	"github.com/google/uuid"
)

type Service struct {
	ExampleRepository exampleRepository
	Logger            *logging.Logger
	// logger etc.
}

func (s Service) CreateExample(ctx context.Context, dto ExampleDTO) (Example, error) {
	// Validate
	if err := dto.Validate(); err != nil {
		return Example{}, WrapError(err, ErrValidation)
	}

	// Create
	ex := newExample(dto)

	// Store
	if err := s.ExampleRepository.Save(ctx, ex); err != nil {
		return Example{}, fmt.Errorf("could not store example %w", err)
	}

	return ex, nil
}

func (s Service) GetExampleByID(ctx context.Context, id uuid.UUID) (Example, error) {
	ex, err := s.ExampleRepository.FindOneByID(ctx, id)
	if err != nil {
		return Example{}, fmt.Errorf("could not get example by id: %s, %w", id, err)
	}

	return ex, nil
}
