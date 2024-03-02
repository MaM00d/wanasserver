package user

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type UserStorage interface {
	InsertUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUserByEmail(string) (*User, error)
}

type userStore struct {
	db *sql.DB
}

func newUserStore(db *sql.DB) *userStore {
	return &userStore{
		db: db,
	}
}

func (s *userStore) InitDb() error {
	s.dropUserTabel()
	return s.createUserTabel()
}

func (s *userStore) dropUserTabel() error {
	query := `
    drop table if exists "User";
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *userStore) createUserTabel() error {
	query := `
        CREATE TABLE IF NOT EXISTS "User" (
            "ID" serial   NOT NULL,
            "Name" char(50)   NOT NULL,
            "Email" char(50)   NOT NULL,
            "Password" char(100)   NOT NULL,
            "Phone" char(50)   NULL,
            "CreatedAt" timestamp   NOT NULL,
            CONSTRAINT "pk_User" PRIMARY KEY (
                "ID"
             )
        );
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *userStore) InsertUser(user *User) error {
	query := `insert into "User" 
    ("Name","Phone","Email","Password","CreatedAt")
    values ($1,$2,$3,$4,$5)
    `
	resp, err := s.db.Query(
		query,
		user.Name,
		user.Phone,
		user.Email,
		user.Password,
		user.CreatedAt,
	)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
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

func scanIntoAccount(rows *sql.Rows) (*User, error) {
	eluser := new(User)
	err := rows.Scan(
		&eluser.ID,
		&eluser.Name,
		&eluser.Phone,
		&eluser.Email,
		&eluser.Password,
		&eluser.CreatedAt)

	return eluser, err
}

func (s *userStore) GetUserByEmail(email string) (*User, error) {
	eluser := new(User)
	err := s.db.QueryRow("select * from users where email = $1", email).Scan(
		&eluser.ID,
		&eluser.Name,
		&eluser.Phone,
		&eluser.Email,
		&eluser.Password,
		&eluser.CreatedAt)

	if err == sql.ErrNoRows {

		slog.Info("no user found with this email", "email", email)
		return nil, fmt.Errorf("user with email [%d] not found", email)
	}
	// for rows.Next() {
	// 	return scanIntoAccount(rows)
	// }

	return eluser, nil
}
