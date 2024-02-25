package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type (
	apiFunc  func(http.ResponseWriter, *http.Request) error
	ApiError struct {
		Error string
	}
)

type APIServer struct {
	listenAddr string
	store      Storage
}

// create object of struct apiserver to set the listen addr
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run the APIServer
func (s *APIServer) Run() {
	// create router
	router := mux.NewRouter()
	// Routes
	router.HandleFunc("/account", makeHTTPhHandleFunc(s.handleAccount))
	// logging
	log.Println("JSON API server runngin on port: ", s.listenAddr)
	// start listening on addresss and sending to router
	http.ListenAndServe(s.listenAddr, router)
}

func makeHTTPhHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleRegister(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed")
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	// id := mux.Vars(r)["id"]
	// fmt.Println(id)
	// return WriteJSON(w, http.StatusOK, &User())
	return nil
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	userReq := new(User)
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		return err
	}

	user := NewUser(
		userReq.FirstName,
		userReq.LastName,
		userReq.Email,
		userReq.Username,
		userReq.Password,
	)
	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
