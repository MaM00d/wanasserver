package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/register", makeHTTPHandleFunc(s.handleUser))
	// logging
	log.Println("JSON API server runngin on port: ", s.listenAddr)
	// start listening on addresss and sending to router
	http.ListenAndServe(s.listenAddr, router)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
