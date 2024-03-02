package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *ElUser) Register(w http.ResponseWriter, r *http.Request) error {
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
