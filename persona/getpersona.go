package persona

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type PersonaRequest struct {
	UserID int `json:"userid"`
}
type PersonaResponse struct {
	Name      string    `json:"personaname"`
	UserID    int       `json:"userid"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func (s *ElPersona) getPersonas(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Persona")
	var req PersonaRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("reading body of persona")
	}
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	resp := PersonaResponse{
		Name:      "Dodda",
		UserID:    2,
		CreatedAt: time.Now().UTC(),
	}

	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
