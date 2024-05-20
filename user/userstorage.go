package user

import (
	"fmt"
	"log/slog"
)

func (s ElUser) createUserTable() error {
	query := `
          CREATE TABLE IF NOT EXISTS users (id serial NOT NULL, name char(50) NOT NULL, email char(50) NOT NULL, password char(100) NOT NULL, phone integer NULL, createdat timestamp NOT NULL, CONSTRAINT pk_user PRIMARY KEY (id));
    `
	err := s.db.Exec(query)
	return err
}

func (s ElUser) dropUserTabel() error {
	query := `
          DROP TABLE IF EXISTS users;
    `
	err := s.db.Exec(query)
	return err
}

func (s ElUser) Insert(user *User) error {
	query := `insert into Users 
          (name,email,password,phone,createdat)
            VALUES ($1,$2,$3,$4,$5)
    `
	return s.db.Query(
		query,
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
	if err != nil {
		return nil, err
	}
	if len(eluser) == 0 {
		slog.Info("no user found with this email", "email", id)
		return nil, fmt.Errorf("user with email [%d] not found", id)
	}
	// for rows.Next() {
	// 	return scanIntoAccount(rows)
	// }

	return eluser[0], nil
}

func (s ElUser) SelectUserByEmail(email string) (*User, error) {
	var eluser []User
	err := s.db.QueryScan(&eluser, `select * from Users where email = $1`, email)
	if err != nil {
		return nil, err
	}
	if len(eluser) == 0 {
		slog.Info("no user found with this email", "email", email)
		return nil, s.db.NotFound
	}

	// for rows.Next() {
	// 	return scanIntoAccount(rows)
	// }

	return &eluser[0], nil
}
