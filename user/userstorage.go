package user

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (s ElUser) InitDb() error {
	if err := s.createUserTabel(); err != nil {
		return err
	}
	return nil
}

func (s ElUser) DropDb() error {
	if err := s.dropUserTabel(); err != nil {
		return err
	}
	return nil
}

func (s ElUser) dropUserTabel() error {
	query := `
    drop table if exists Users;
    `
	err := s.db.Exec(query)
	return err
}

func (s ElUser) createUserTabel() error {
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

func (s ElUser) InsertUser(user *User) error {
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

func (s ElUser) DeleteUser(int) error {
	return nil
}

func (s ElUser) UpdateUser(*User) error {
	return nil
}

// func (s ElUser) GetUserById(int) (*User, error) {
// 	return nil, nil
// }

func (s ElUser) SelectUserById(id int) (*User, error) {
	var eluser []*User
	err := s.db.QueryScan(&eluser, `select * from Users where ID = $1`, id)
	if err == sql.ErrNoRows {

		slog.Info("no user found with this email", "email", id)
		return nil, fmt.Errorf("user with email [%d] not found", id)
	}
	// for rows.Next() {
	// 	return scanIntoAccount(rows)
	// }

	return eluser[0], nil
}

func (s ElUser) SelectUserByEmail(email string) (*User, error) {
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
