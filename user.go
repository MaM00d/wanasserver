package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
type UserView struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func (a *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(pw)) == nil
}

func NewUser(firstName, lastName, email, username, password string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(encpw),
		Username:  username,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUser(w, r)
	}
	if r.Method == "POST" {
		return s.handleRegister(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUser(w, r)
	}

	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	userReq := new(User)
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		return err
	}

	// create user object from user struct
	user, err := NewUser(
		userReq.FirstName,
		userReq.LastName,
		userReq.Email,
		userReq.Username,
		userReq.Password,
	)
	if err != nil {
		return err
	}

	// save user to database
	if err := s.store.InsertUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func createJWT(eluser *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":    15000,
		"accountEmail": eluser.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		return err
	}

	acc, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token: token,
		Email: acc.Email,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
