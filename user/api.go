package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"

	api "Server/elapi"
)

type ElUser struct {
	ap    *api.ElApi
	store UserStorage
}

// create object of struct apiserver to set the listen addr
func NewElUser(db *sql.DB, elapi *api.ElApi) *ElUser {
	store := newUserStore(db)
	eluser := &ElUser{
		ap:    elapi,
		store: store,
	}
	eluser.addroutes()
	return eluser
}

func (eluser *ElUser) addroutes() {
	eluser.ap.Route("/login", eluser.HandleLogin)
	eluser.ap.Route("/register", eluser.HandleUser)
	eluser.ap.Route("/user/{email}", eluser.JWTAuthMiddleware)
}

func (s *ElUser) HandleUser(w http.ResponseWriter, r *http.Request) error {
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

func (s *ElUser) handleRegister(w http.ResponseWriter, r *http.Request) error {
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
	return s.ap.WriteJSON(w, http.StatusOK, user)
}

func (s *ElUser) handleGetUserByEmailFromVars(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		email := s.getEmailFromVars(r)

		account, err := s.store.GetUserByEmail(email)
		if err != nil {
			return err
		}

		return s.ap.WriteJSON(w, http.StatusOK, account)
	}

	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (eluser *ElUser) JWTAuthMiddleware(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("calling JWT auth middleware")
	tokenString := r.Header.Get("x-jwt-token")

	token, err := detokenizejwt(tokenString)
	if err != nil {
		eluser.ap.PermissionDenied(w)
	}
	if !token.Valid {

		fmt.Println("err 2")
		eluser.ap.PermissionDenied(w)
	}
	userEmail := eluser.getEmailFromVars(r)

	user, err := eluser.store.GetUserByEmail(userEmail)
	if err != nil {

		fmt.Println("err 3")
		eluser.ap.PermissionDenied(w)
		eluser.ap.PermissionDenied(w)
	}

	claims := token.Claims.(jwt.MapClaims)

	if user.Email != claims["userEmail"] {

		fmt.Println("err 4")
		eluser.ap.PermissionDenied(w)
	}

	if err != nil {
		eluser.ap.WriteJSON(w, http.StatusForbidden, api.ApiError{Error: "invalid token"})
	}
	return eluser.handleGetUserByEmailFromVars(w, r)
}

func (s *ElUser) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ElUser) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
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

func (eluser *ElUser) getEmailFromVars(r *http.Request) string {
	elemail := eluser.ap.GetFromVars(r, "email")
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
