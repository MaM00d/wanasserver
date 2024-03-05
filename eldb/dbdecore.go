package eldb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Storage struct {
	// db *sql.DB
	db  *pgxpool.Pool
	ctx context.Context
}

func NewPostgresStore() (*Storage, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db:  db,
		ctx: ctx,
	}, nil

	// connStr := "user=tme password='1598753' dbname=wanas sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := db.Ping(); err != nil {
	// 	return nil, err
	// }
	// return &Storage{
	// 	db: db,
	// }, nil
}

func (s *Storage) Exec(query string) error {
	resp, err := s.db.Exec(context.Background(), query)
	slog.Info("SQL", "Response", resp)
	return err
}

func (s *Storage) Query(query string, args ...any) error {
	_, err := s.db.Query(s.ctx, query, args...)
	if err != nil {
		slog.Error("SQL", "Query", err)
	}
	return err
}

func (s *Storage) QueryRow(query string, args ...any) pgx.Row {
	resp := s.db.QueryRow(s.ctx, query, args)

	return resp
}

func (s *Storage) Scan(rows *sql.Rows, obj ...any) error {
	err := rows.Scan(obj)
	return err
}

func (s *Storage) QueryScan(obj interface{}, query string, args ...any) error {
	err := pgxscan.Select(s.ctx, s.db, obj, query, args...)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("not found")
	}
	if err != nil {
		slog.Error("SQL", "QueryScan", err)
		return err
	}
	return nil
}
