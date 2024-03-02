package user

import (
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"

	api "Server/elapi"
)

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
	userEmail := eluser.getIdFromVars(r)

	user, err := eluser.store.GetUserById(userEmail)
	if err != nil {

		fmt.Println("err 3")
		eluser.ap.PermissionDenied(w)
		eluser.ap.PermissionDenied(w)
	}

	claims := token.Claims.(jwt.MapClaims)

	if user.Email != claims["userid"] {

		fmt.Println("err 4")
		eluser.ap.PermissionDenied(w)
	}

	if err != nil {
		eluser.ap.WriteJSON(w, http.StatusForbidden, api.ApiError{Error: "invalid token"})
	}
	return eluser.getUserByEmailFromVars(w, r)
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

	account, err := s.store.GetUserById(id)
	if err != nil {
		return err
	}

	return s.ap.WriteJSON(w, http.StatusOK, account)

	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }
}

func (eluser *ElUser) getIdFromVars(r *http.Request) string {
	elemail := eluser.ap.GetFromVars(r, "id")
	return elemail
}
