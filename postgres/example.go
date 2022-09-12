package postgres

import (
	"context"

	"github.com/bratteby/go-service-template/example"
	"github.com/google/uuid"
)

type ExampleRepository struct {
	DB pool
}

func (r *ExampleRepository) FindOneByID(ctx context.Context, id uuid.UUID) (example.Example, error) {
	query := `
		SELECT id, name
		FROM example
		WHERE id = $1
	`

	var ex example.Example
	if err := r.DB.QueryRow(ctx, query, id).Scan(&ex.ID, &ex.Name); err != nil {
		return example.Example{}, wrapPgxError(err)
	}

	return ex, nil
}

func (r *ExampleRepository) Save(ctx context.Context, ex example.Example) error {
	sql := `
		INSERT INTO example(id, name) values (
			$1, $2
		) 
	`

	_, err := r.DB.Exec(ctx, sql, ex.ID, ex.Name)
	if err != nil {
		return wrapPgxError(err)
	}

	return nil
}
