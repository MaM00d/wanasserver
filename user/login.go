package user

import (
	"encoding/json"
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
	req := new(LoginRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	slog.Info("Login", "decoded", req.Email)
	acc, err := s.SelectUserByEmail(req.Email)
	if err == s.db.NotFound {
		slog.Error("no user found with this email")
		return s.ap.WriteJSON(w, http.StatusNotFound, "No user found with this Email")
	}

	slog.Info("Login", "selected", acc.Email)
	slog.Info(req.Password)
	if !acc.ValidPassword(req.Password) {
		return s.ap.WriteError(w, http.StatusUnauthorized, "Wrong Password")
	}

	token, err := tokenizejwt(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token: token,
		Email: acc.Email,
	}
	slog.Info("Login", "Success", resp.Email)
	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
