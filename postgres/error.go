package postgres

import (
	"net"

	"errors"

	"github.com/bratteby/go-service-template/example"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

func wrapPgxError(err error) error {
	if err == pgx.ErrNoRows {
		return example.WrapError(err, example.ErrNotFound)
	}

	var opError *net.OpError
	if errors.As(err, &opError) {
		if opError.Op == "dial" {
			return example.WrapError(err, example.ErrTemporary)
		}
	}

	var pgxError *pgconn.PgError
	if errors.As(err, &pgxError) {
		code := pgxError.Code

		switch {
		case pgerrcode.IsConnectionException(code), pgerrcode.IsConnectionException(code):
			return example.WrapError(err, example.ErrTemporary)
		}
	}

	return err
}
