package chat

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type CreateChatRequest struct {
	PersonaID string `db:"personaid" json:"personaid"`
}

func (s *ElChat) createchat(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	slog.Info("Handling Create Chat")
	chatReq := new(CreateChatRequest)
	if err := json.NewDecoder(r.Body).Decode(chatReq); err != nil {
		slog.Error("decoding request body", "Model", "Chat")
		return err
	}
	elpersonaid, err := strconv.Atoi(chatReq.PersonaID)
	if err != nil {
		return err
	}

	chat, err := NewChat(
		elpersonaid,
	)
	if err != nil {
		return err
	}

	if err := s.store.InsertChat(chat); err != nil {
		return err
	}

	slog.Info("Successfully Registered")
	return s.ap.WriteJSON(w, http.StatusOK, chat)
}
