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
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}

	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	if err != nil {
		return err
	}

	chat := NewChat(
		elpersonaid,
		eluserid,
	)

	if err := s.InsertChat(chat); err != nil {
		return err
	}

	slog.Info("created chat")
	return s.ap.WriteJSON(w, http.StatusOK, chat)
}
