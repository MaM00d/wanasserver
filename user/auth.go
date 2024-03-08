package user

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"

	api "Server/elapi"
)

func (eluser *ElUser) JWTAuthMiddleware(
	elfunc api.ApiFunc,
) api.ApiFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) error {
		slog.Info("calling JWT auth middleware")
		tokenString := r.Header.Get("x-jwt-token")

		token, err := detokenizejwt(tokenString)
		if err != nil {
			if len(tokenString) == 0 {
				slog.Error("no token in header")
			} else {
				slog.Error("detokenize", err)
			}

			return errors.New("invalid token")
		}
		if !token.Valid {

			slog.Error("error validate")
			return errors.New("invalid token")
		}
		userid := eluser.getIdFromVars(r)

		euser, err := eluser.SelectUserById(userid)
		if err != nil {

			slog.Error("error getting user by id")
			return errors.New("invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)

		if euser.ID != int(claims["userid"].(float64)) {
			slog.Error("id in token not equal id in route")
			return errors.New("invalid token")
		}
		if err != nil {
			return errors.New("invalid token")
		}
		return elfunc(w, r)
	}
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

func (s *ElUser) getUserByEmailFromVars(w http.ResponseWriter, r *http.Request) error {
	id := s.getIdFromVars(r)

	account, err := s.SelectUserById(id)
	if err != nil {
		return err
	}

	return s.ap.WriteJSON(w, http.StatusOK, account)

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
