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

	// save user to database
	if err := s.InsertUser(user); err != nil {
		return err
	}

	slog.Info("Successfully Registered")
	return s.ap.WriteJSON(w, http.StatusOK, user)
}
