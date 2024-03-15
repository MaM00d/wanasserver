package msg

import (
	"time"
)

type Msg struct {
	ID        int       `db:"id"        json:"id"`
	ChatID    int       `db:"chatid"    json:"chatid"`
	PersonaID int       `db:"personaid"    json:"personaid"`
	UserID    int       `db:"userid"    json:"userid"`
	Message   string    `db:"message" json:"message"`
	CreatedAt time.Time `db:"createdat" json:"createdat"`
}

type MsgView struct {
	Message   string    `db:"message" json:"message"`
	CreatedAt time.Time `db:"createdat" json:"createdat"`
}

func NewMsg(chatid, personaid, userid int, message string) (*Msg, error) {
	return &Msg{
		ChatID:    chatid,
		PersonaID: personaid,
		UserID:    userid,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}, nil
}
