package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUserById(int) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=tme password='1598753' dbname=wanas sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) InitDb() error {
	return s.createUserTabel()
}

func (s *PostgresStore) createUserTabel() error {
	query := `create table if not exists users (
        id serial primary key,
        first_name varchar(50),
        last_name varchar(50),
        username varchar(50),
        email varchar(50),
        password varchar(50),
        created_at timestamp
    )`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(user *User) error {
	query := `insert into users 
    (first_name,last_name,username,email,password,created_at)
values ($1,$2,$3,$4,$5,$6)
    `
	resp, err := s.db.Query(
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStore) DeleteUser(int) error {
	return nil
}

func (s *PostgresStore) UpdateUser(*User) error {
	return nil
}

func (s *PostgresStore) GetUserById(int) (*User, error) {
	return nil, nil
}
