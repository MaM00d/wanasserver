package user

import (
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id"        json:"id"`
	Name      string    `db:"name"      json:"name"`
	Phone     int       `db:"phone"     json:"phone"`
	Email     string    `db:"email"     json:"email"`
	Password  string    `db:"password"  json:"password"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

type UserView struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     int       `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUser(name, email, password string, phone int) (*User, error) {
	slog.Info(password)
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("error encrypting password")
		return nil, err
	}
	return &User{
		Name:      name,
		Email:     email,
		Password:  string(encpw),
		Phone:     phone,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (a *User) ValidPassword(pw string) bool {
	slog.Info("mamaaaa: ", "pass", a.Password)
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(pw)) == nil
}
