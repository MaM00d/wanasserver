package chat

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type CreateChatRequest struct {
	Name      string    `db:"name"      json:"name"`
	UserID    string    `db:"userid"    json:"userid"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

func (s *ElChat) createchat(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	slog.Info("Handling Create Chat")
	chatReq := new(CreateChatRequest)
	if err := json.NewDecoder(r.Body).Decode(chatReq); err != nil {
		slog.Error("decoding request body", "Model", "Chat")
		return err
	}
	eluserid, err := strconv.Atoi(chatReq.UserID)
	if err != nil {
		return err
	}

	chat, err := NewChat(
		chatReq.Name,
		eluserid,
	)
	if err != nil {
		return err
	}

	if err := s.InsertChat(chat); err != nil {
		return err
	}

	slog.Info("Successfully Registered")
	return s.ap.WriteJSON(w, http.StatusOK, chat)
}
