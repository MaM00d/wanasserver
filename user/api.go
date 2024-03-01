package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// Run the APIServer
func (s *APIServer) Run() {
	// create router
	router := mux.NewRouter()
	// Routes
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/register", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc(
		"/user/{email}",
		JWTAuthMiddleware(makeHTTPHandleFunc(s.handleGetUserByEmailFromVars), s.store),
	)
	// logging
	slog.Info("JSON API server runngin", "PORT", s.listenAddr)
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

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return s.Login(w, r)
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
	slog.Info("Handling Register")
	userReq := new(User)
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		slog.Error("decoding request body")
		return err
	}

	// create user object from user struct
	user, err := NewUser(
		userReq.Name,
		userReq.Email,
		userReq.Password,
		userReq.Phone,
	)
	if err != nil {
		return err
	}

	// save user to database
	if err := s.store.InsertUser(user); err != nil {
		return err
	}

	slog.Info("Successfully Registered")
	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) Login(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")
	var req LoginRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	acc, err := s.store.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := tokenizejwt(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token: token,
		Email: acc.Email,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleGetUserByEmailFromVars(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		email := getEmailFromVars(r)

		account, err := s.store.GetUserByEmail(email)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }

	return fmt.Errorf("method not allowed %s", r.Method)
}

func JWTAuthMiddleware(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")
		tokenString := r.Header.Get("x-jwt-token")

		token, err := detokenizejwt(tokenString)
		if err != nil {
			fmt.Println("err 1")
			permissionDenied(w)
			return
		}
		if !token.Valid {

			fmt.Println("err 2")
			permissionDenied(w)
			return
		}
		userEmail := getEmailFromVars(r)

		user, err := s.GetUserByEmail(userEmail)
		if err != nil {

			fmt.Println("err 3")
			permissionDenied(w)
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if user.Email != claims["userEmail"] {

			fmt.Println("err 4")
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func tokenizejwt(eluser *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userEmail": eluser.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func getEmailFromVars(r *http.Request) string {
	elemail := mux.Vars(r)["email"]
	return elemail
}

func detokenizejwt(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
