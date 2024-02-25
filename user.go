package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	// ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
type UserView struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUser(firstName, lastName, email, username, password string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Username:  username,
		CreatedAt: time.Now().UTC(),
	}
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleLogin(w, r)
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
	user := NewUser(
		userReq.FirstName,
		userReq.LastName,
		userReq.Email,
		userReq.Username,
		userReq.Password,
	)
	// save user to database
	if err := s.store.InsertUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
