package chat

import (
	"time"
)

type Chat struct {
	ID        int       `db:"id"        json:"id"`
	PersonaID int       `db:"personaid"    json:"personaid"`
	UserID    int       `db:"userid"    json:"userid"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

type ChatView struct {
	ID        int       `db:"id"        json:"id"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

func NewChat(personaid, userid int) (*Chat, error) {
	return &Chat{
		PersonaID: personaid,
		UserID:    userid,
		CreatedAt: time.Now().UTC(),
	}, nil
}
