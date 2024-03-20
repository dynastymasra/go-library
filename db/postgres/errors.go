package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// IsUniqueViolation checks if the provided error is a unique violation error in PostgresSQL.
// It returns true if the error is a unique violation error (error code "23505"), and false otherwise.
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// IsForeignKeyViolation checks if the provided error is a foreign key violation error in PostgresSQL.
// It returns true if the error is a foreign key violation error (error code "23503"), and false otherwise.
func IsForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

// IsInvalidTextRepresentation checks if the provided error is an invalid text representation error in PostgresSQL.
// It returns true if the error is an invalid text representation error (error code "22P02"), and false otherwise.
func IsInvalidTextRepresentation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "22P02"
}

// IsNotNullViolation checks if the provided error is a not null violation error in PostgresSQL.
// It returns true if the error is a not null violation error (error code "23502"), and false otherwise.
func IsNotNullViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23502"
}
