package example

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetExample(t *testing.T) {
	// Arrange
	existingID := uuid.New()
	exampleExample := Example{
		ID:   existingID,
		Name: "Test",
	}

	repo := &exampleRepositoryMock{
		FindOneByIDFunc: func(ctx context.Context, id uuid.UUID) (Example, error) {
			switch id {
			case existingID:
				return exampleExample, nil
			default:
				return Example{}, ErrNotFound
			}

		},
		SaveFunc: func(ctx context.Context, ex Example) error {
			return nil
		},
	}

	s := Service{
		ExampleRepository: repo,
	}

	tests := []struct {
		name          string
		givenID       uuid.UUID
		expected      Example
		expectedError error
	}{
		{
			name:          "should return error on non existing example",
			givenID:       uuid.New(),
			expected:      Example{},
			expectedError: ErrNotFound,
		},
		{
			name:          "should return expected example",
			givenID:       existingID,
			expected:      exampleExample,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, err := s.GetExampleByID(context.Background(), tt.givenID)

			// Assert
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}

}
