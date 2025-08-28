package helpers

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeUniqueViolation     = "23505"
	CodeForeignKeyViolation = "23503"
	CodeNotNullViolation    = "23502"
)

func Code(err error) string {
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		return pg.Code
	}
	return ""
}

func Constraint(err error) string {
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		return pg.ConstraintName
	}
	return ""
}

func IsUnique(err error) bool     { return Code(err) == CodeUniqueViolation }
func IsForeignKey(err error) bool { return Code(err) == CodeForeignKeyViolation }
func IsNotNull(err error) bool    { return Code(err) == CodeNotNullViolation }

func IsOnConstraint(err error, nameFrag string) bool {
	if !IsUnique(err) {
		return false
	}
	return strings.Contains(Constraint(err), nameFrag)
}
