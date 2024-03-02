package msg

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type MsgRequest struct {
	text      string `json:"text"`
	email     string `json:"email"`
	personaid int    `json:"persona"`
	CreatedAt string `json:"CreatedAt"`
}
type MsgResponse struct {
	text      string    `json:"text"`
	email     string    `json:"email"`
	personaid int       `json:"persona"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func (s *ElChat) Send(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Msg")
	var req MsgRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("reading body of msg")
	}
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	resp := MsgResponse{
		text:      "Hello from Ai",
		email:     "email",
		personaid: 1,
		CreatedAt: time.Now().UTC(),
	}

	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
