package persona

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type CreatePersonaRequest struct {
	Name      string    `db:"name"      json:"name"`
	UserID    string    `db:"userid"    json:"userid"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

func (s *ElPersona) createpersona(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	slog.Info("Handling Create Persona")
	personaReq := new(CreatePersonaRequest)
	if err := json.NewDecoder(r.Body).Decode(personaReq); err != nil {
		slog.Error("decoding request body", "Model", "Persona")
		return err
	}
	eluserid, err := strconv.Atoi(personaReq.UserID)
	if err != nil {
		return err
	}

	persona, err := NewPersona(
		personaReq.Name,
		eluserid,
	)
	if err != nil {
		return err
	}

	if err := s.InsertPersona(persona); err != nil {
		return err
	}

	slog.Info("Successfully Registered")
	return s.ap.WriteJSON(w, http.StatusOK, persona)
}
