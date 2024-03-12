package persona

import (
	"time"
)

type Persona struct {
	ID        int       `db:"id"        json:"id"`
	Name      string    `db:"name"      json:"name"`
	UserID    int       `db:"userid"    json:"userid"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

type PersonaView struct {
	ID        int       `db:"id"        json:"id"`
	Name      string    `db:"name"      json:"name"`
	UserID    int       `db:"userid"    json:"userid"`
	CreatedAt time.Time `db:"createdat" json:"createdAt"`
}

func NewPersona(name string, userid int) (*Persona, error) {
	return &Persona{
		Name:      name,
		UserID:    userid,
		CreatedAt: time.Now().UTC(),
	}, nil
}
