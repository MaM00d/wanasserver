package eldb

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Storage struct {
	// db *sql.DB
	db       *pgxpool.Pool
	ctx      context.Context
	entities *[]entity
	NotFound error
}

type entity interface {
	create() error
	drop() error
}

func NewPostgresStore() (*Storage, error) {
	slog.Info("connecting to database")
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	db, err := pgxpool.ConnectConfig(ctx, config)
	eltime := 0
	for db == nil {
		time.Sleep(2 * time.Second)
		db, err = pgxpool.ConnectConfig(ctx, config)
		eltime = eltime + 2

		if err != nil && eltime == 10 {

			slog.Error("Time Out Can't connet to database")
			return nil, err
		}
		slog.Warn("waiting to connect to database")
	}
	slog.Info("connected to database successfully")

	return &Storage{
		db:       db,
		ctx:      ctx,
		NotFound: errors.New("No Row Found in Db"),
	}, nil
}

func (s *Storage) Exec(query string) error {
	resp, err := s.db.Exec(context.Background(), query)
	slog.Info("SQL", "Response", resp)
	if err != nil {
		slog.Error("SQL", "query", query, "Exec", err)
	}
	return err
}

func (s *Storage) Query(query string, args ...any) error {
	rows, err := s.db.Query(s.ctx, query, args...)
	if err != nil {
		slog.Error("SQL", "Query", err)
		return err
	}
	defer rows.Close()
	return nil
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
	if err != nil {
		return err
	}

	return nil
}
