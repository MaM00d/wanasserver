package chat

import (
	"Server/user"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *ElChat) createchat(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	slog.Info("Handling Create Chat")
	eluserid := user.Getidfromheader(r)

	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	chat, err := NewChat(
		elpersonaid,
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
