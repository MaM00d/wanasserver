package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (s *ElUser) Login(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")
	var req LoginRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	acc, err := s.SelectUserByEmail(req.Email)
	if err != nil {
		return err
	}

	slog.Info(req.Password)
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

	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
