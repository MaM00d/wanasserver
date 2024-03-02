package main

import (
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func NewPostgresStore() (*Storage, error) {
	connStr := "user=tme password='1598753' dbname=wanas sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}
