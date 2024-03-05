package user

import (
	"database/sql"
	"fmt"
	"log/slog"

	db "Server/eldb"
)

type UserStorage interface {
	InsertUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUserById(int) (*User, error)
	GetUserByEmail(string) (*User, error)
}

type userStore struct {
	db db.Storage
}

func newUserStore(db db.Storage) *userStore {
	elstore := &userStore{
		db: db,
	}
	elstore.InitDb()
	return elstore
}

func (s *userStore) InitDb() error {
	s.dropUserTabel()
	return s.createUserTabel()
}

func (s *userStore) dropUserTabel() error {
	query := `
    drop table if exists Users;
    `
	err := s.db.Exec(query)
	return err
}

func (s *userStore) createUserTabel() error {
	query := `
        CREATE TABLE IF NOT EXISTS Users (
            ID serial   NOT NULL,
            Name char(50)   NOT NULL,
            Email char(50)   NOT NULL,
            Password char(100)   NOT NULL,
            Phone integer   NULL,
            CreatedAt timestamp   NOT NULL,
            CONSTRAINT pk_User PRIMARY KEY (
                ID
             )
        );
    `
	err := s.db.Exec(query)
	return err
}

func (s *userStore) InsertUser(user *User) error {
	query := `insert into Users 
    (id,Name,Email,Password,Phone,CreatedAt)
    values ($1,$2,$3,$4,$5,$6)
    `
	return s.db.Query(
		query,
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.CreatedAt,
	)
}

func (s *userStore) DeleteUser(int) error {
	return nil
}

func (s *userStore) UpdateUser(*User) error {
	return nil
}

// func (s *userStore) GetUserById(int) (*User, error) {
// 	return nil, nil
// }

func (s *userStore) SelectUserById(id int) (*User, error) {
	eluser := new(User)
	err := s.db.QueryScan(eluser, `select * from Users where ID = $1`, id)
	if err == sql.ErrNoRows {

		slog.Info("no user found with this email", "email", id)
		return nil, fmt.Errorf("user with email [%d] not found", id)
	}
	// for rows.Next() {
	// 	return scanIntoAccount(rows)
	// }

	return eluser, nil
}

func (s *userStore) SelectUserByEmail(email string) (*User, error) {
	var eluser []*User
	err := s.db.QueryScan(&eluser, `select * from Users where email = $1`, email)
	if err == fmt.Errorf("not found") {
		return nil, fmt.Errorf("user with email [%d] not found", email)
	}
	// for rows.Next() {
	// 	return scanIntoAccount(rows)
	// }

	return eluser[0], nil
}
