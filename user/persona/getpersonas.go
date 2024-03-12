package persona

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type GetPersonasRequest struct {
	id int `json:"userid"`
}
type GetPersonasResponse struct {
	personas []Persona `json:"personas"`
}

func (s *ElPersona) getpersonas(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")
	var req GetPersonasRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	elpersonas, err := s.GetPersonasByUserId(req.id)
	if err != nil {
		return err
	}

	resp := GetPersonasResponse{
		personas: elpersonas,
	}

	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
