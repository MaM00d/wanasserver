package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *ElUser) Register(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	slog.Info("Handling Register")
	userReq := new(RegisterRequest)
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		slog.Error("decoding request body")
		return err
	}
	elphone, err := strconv.Atoi(userReq.Phone)
	if err != nil {
		return err
	}

	// create user object from user struct
	user := NewUser(
		userReq.Name,
		userReq.Email,
		userReq.Password,
		elphone,
	)
	if err := user.Encrippass(); err != nil {
		return err
	}

	if usr, _ := s.SelectUserByEmail(userReq.Email); usr != nil {
		slog.Error("user with email exist")
		return s.ap.WriteError(w, http.StatusConflict, "User exsists with that email")
	}
	// save user to database
	if err := s.Insert(user); err != nil {
		return err
	}

	slog.Info("Successfully Registered")

	token, err := tokenizejwt(user)
	if err != nil {
		return err
	}
	resp := LoginResponse{
		Token: token,
		Email: user.Email,
	}
	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
