package msg

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type ChatStorage interface {
	InsertMsg(*Msg) error
	DeleteMsg(int) error
	UpdateMsg(*Msg) error
}

type chatStore struct {
	db *sql.DB
}

func newChatStore(db *sql.DB) *chatStore {
	return &chatStore{
		db: db,
	}
}

func (s *chatStore) InitDb() error {
	s.dropChatTabel()
	return s.createMsgTabel()
}

func (s *chatStore) dropChatTabel() error {
	query := `
    drop table if exists "Chat";
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *chatStore) createMsgTabel() error {
	query := `
CREATE TABLE "Chat" (
    "ChatID" serial   NOT NULL,
    "Message" char(50)   NOT NULL,
    "date" timestamp   NOT NULL,
    "PersonaID" int   NOT NULL,
    CONSTRAINT "pk_Chat" PRIMARY KEY (
        "ChatID"
     )
);
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *chatStore) InsertMsg(chat *Msg) error {
	query := `insert into "Chat" 
    ("Message","date",PersonaID)
    values ($1,$4,$5)
    `
	resp, err := s.db.Query(
		query,
		chat.text,
		chat.CreatedAt,
		chat.personaid,
	)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *chatStore) DeleteMsg(int) error {
	return nil
}

func (s *chatStore) UpdateMsg(*Msg) error {
	return nil
}

// func (s *chatStore) GetMsgById(int) (*Msg, error) {
// 	return nil, nil
// }

func scanIntoAccount(rows *sql.Rows) (*Msg, error) {
	elchat := new(Msg)
	err := rows.Scan(
		&elchat.text,
		&elchat.personaid,
		&elchat.CreatedAt)

	return elchat, err
}

func (s *chatStore) GetChat(email, persona string) (*Msg, error) {
	elchat := new(Msg)
	chats, err := s.db.Query("select * from chats where email = $1 and persona=$2", email, persona)

	if err == sql.ErrNoRows {

		slog.Info("no chat found with this email", "email", email)
		return nil, fmt.Errorf("chat with email [%d] not found", email)
	}

	for chats.Next() {
		return scanIntoAccount(chats)
	}

	return elchat, nil
}
