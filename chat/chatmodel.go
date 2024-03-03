package chat

import (
	"time"
)

type Chat struct {
	ID        int       `db:"id"        json:"id"`
	PersonaID int       `db:"personaid" json:"personaid"`
	CreatedAt time.Time `db:"createdat" json:"createdat"`
}

type Msg struct {
	ID        int       `db:"id"      json:"id"`
	Message   string    `db:"message" json:"message"`
	ChatID    int       `db:"chatid"  json:"chatid"`
	CreatedAt time.Time `db:"date"    json:"createdat"`
}

func NewChat(personaid int) (*Chat, error) {
	return &Chat{
		PersonaID: personaid,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func NewMsg(message string, chatid int) (*Msg, error) {
	return &Msg{
		Message:   message,
		ChatID:    chatid,
		CreatedAt: time.Now().UTC(),
	}, nil
}
