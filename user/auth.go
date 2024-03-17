package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
)

func (eluser *ElUser) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Middleware")
		err := eluser.Authenticate(w, r)
		if err != nil {
			slog.Info("notAuthenticated", "Error", err)
			if r.URL.Path == "/login" || r.URL.Path == "/register" {

				slog.Info("NewUser")
				next.ServeHTTP(w, r)
			} else {
				slog.Info("forbidden")
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		} else {
			slog.Info("Authenticated")
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)

		}
	})
}

func Getidfromheader(r *http.Request) int {
	tokenString := r.Header.Get("x-jwt-token")

	token, err := detokenizejwt(tokenString)
	if err != nil {
		return -1
	}

	claims := token.Claims.(jwt.MapClaims)
	return int(claims["userid"].(float64))
}

func (eluser *ElUser) Authenticate(
	w http.ResponseWriter,
	r *http.Request,
) error {
	slog.Info("Authenticating")
	tokenString := r.Header.Get("x-jwt-token")

	token, err := detokenizejwt(tokenString)
	if err != nil {
		if len(tokenString) == 0 {
			slog.Error("no token in header")
		} else {
			slog.Error("detokenize", err)
		}

		return eluser.ap.WriteError(w, http.StatusUnauthorized, "Invalid session token")
	}
	if !token.Valid {

		slog.Error("error validate")
		return eluser.ap.WriteError(w, http.StatusUnauthorized, "Invalid session token")
	}

	claims := token.Claims.(jwt.MapClaims)
	if _, err := eluser.SelectUserById(int(claims["userid"].(float64))); err != nil {

		slog.Error("no user in database with that id")
		return eluser.ap.WriteError(w, http.StatusUnauthorized, "Invalid session token")
	}

	if err != nil {
		return eluser.ap.WriteError(w, http.StatusUnauthorized, "Invalid session token")
	}
	return nil
}

func tokenizejwt(eluser *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userid":    eluser.ID,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
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

func (s *ElUser) getUserByIDFromVars(w http.ResponseWriter, r *http.Request) error {
	// id := s.getIdFromVars(r)
	id := Getidfromheader(r)

	account, err := s.SelectUserById(id)
	if err != nil {
		return err
	}

	return s.ap.WriteJSON(w, http.StatusOK, account.ID)

	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }
}

func (eluser *ElUser) getIdFromVars(r *http.Request) int {
	elid := eluser.ap.GetFromVars(r, "id")
	eli, err := strconv.Atoi(elid)
	if err != nil {
		return 0
	}
	return eli
}
