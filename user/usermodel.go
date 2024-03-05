package user

import (
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id"        json:"id"`
	Name      string    `db:"name"      json:"name"`
	Email     string    `db:"email"     json:"email"`
	Password  string    `db:"password"  json:"password"`
	Phone     int       `db:"phone"     json:"phone"`
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

func NewUser(name, email, password string, phone int) *User {
	slog.Info(password)
	elusr := &User{
		ID:        0,
		Name:      name,
		Email:     email,
		Password:  password,
		Phone:     phone,
		CreatedAt: time.Now().UTC(),
	}
	return elusr
}

func (a *User) ValidPassword(pw string) bool {
	slog.Info("mamaaaa: ", "pass", a.Password)
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(pw)) == nil
}

func (a *User) Encrippass() error {
	encpw, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	a.Password = string(encpw)
	if err != nil {
		slog.Error("error encrypting password")
		return err
	}
	return nil
}
