package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-rest-chi/internal/dbgen"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQL struct {
	DB *sql.DB
	Q  *dbgen.Queries
}

type Options struct {
	Driver      string
	DSN         string
	MaxOpen     int
	MaxIdle     int
	MaxIdleTime time.Duration
}

func Open(opts Options) (*SQL, error) {
	if opts.Driver == "" {
		return nil, errors.New("db.Open: empty driver")
	}

	if opts.DSN == "" {
		return nil, errors.New("db.Open: empty DSN")
	}

	db, err := sql.Open(opts.Driver, opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("db.Open: %w", err)
	}

	if opts.MaxOpen > 0 {
		db.SetMaxOpenConns(opts.MaxOpen)
	}
	if opts.MaxIdle > 0 {
		db.SetMaxIdleConns(opts.MaxIdle)
	}
	if opts.MaxIdleTime > 0 {
		db.SetConnMaxIdleTime(opts.MaxIdleTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db.PingContext %w", err)
	}

	return &SQL{DB: db, Q: dbgen.New(db)}, nil
}

func (s *SQL) Check(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}

	var one int
	if err := s.DB.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("db select 1: %w", err)
	}

	if one != 1 {
		return fmt.Errorf("db select 1: unexpected result %v", one)
	}

	return nil
}

func (s *SQL) Version(ctx context.Context) (string, error) {
	var v string

	if err := s.DB.QueryRowContext(ctx, "SELECT VERSION()").Scan(&v); err != nil {
		return "", fmt.Errorf("db version: %w", err)
	}

	return v, nil
}

func (s *SQL) Close() error {
	return s.DB.Close()
}
