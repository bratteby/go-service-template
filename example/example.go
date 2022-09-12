package example

import (
	"fmt"

	"github.com/google/uuid"
)

type Example struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func newExample(dto ExampleDTO) Example {
	return Example{
		ID:   uuid.New(),
		Name: dto.Name,
	}
}

type ExampleDTO struct {
	Name string `json:"name"`
}

func (dto ExampleDTO) Validate() error {
	if dto.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	return nil
}
