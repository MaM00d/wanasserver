package chat

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type GetChatsRequest struct {
	personaid int `json:"userid"`
	chatid    int `json:"userid"`
}
type GetChatsResponse struct {
	chats []Chat `json:"chats"`
}

func (s *ElChat) getchats(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")
	var req GetChatsRequest

	bodybytes, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(bodybytes, &req); err != nil {
		slog.Error("decoding request body")
		return err
	}

	elchats, err := s.store.GetChatsByPersonaId(req.personaid)
	if err != nil {
		return err
	}

	resp := GetChatsResponse{
		chats: elchats,
	}

	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
