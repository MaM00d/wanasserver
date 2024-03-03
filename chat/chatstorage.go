package chat

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type ChatStorage interface {
	InsertChat(*Chat) error
	DeleteChat(int) error
	UpdateChat(*Chat) error
	GetChatById(int) (*Chat, error)
	GetChatByEmail(string) (*Chat, error)
}

type chatStore struct {
	db *sql.DB
}

func newChatStore(db *sql.DB) *chatStore {
	elstore := &chatStore{
		db: db,
	}
	elstore.InitDb()
	return elstore
}

func (s *chatStore) InitDb() error {
	if err := s.dropChatTable(); err != nil {
		return err
	}
	if err := s.dropMsgTable(); err != nil {
		return err
	}
	if err := s.createChatTable(); err != nil {
		return err
	}
	if err := s.createMsgTable(); err != nil {
		return err
	}
	return nil
}

func (s *chatStore) dropChatTable() error {
	query := `
    drop table if exists Chat;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *chatStore) dropMsgTable() error {
	query := `
    drop table if exists msg;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *chatStore) createChatTable() error {
	query := `
            CREATE TABLE chat (
                id serial   NOT NULL,
                personaid int   NOT NULL,
                createdat timestamp   NOT NULL,
                CONSTRAINT pk_chat PRIMARY KEY (
                    id,personaid
                 )
            );
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *chatStore) createMsgTable() error {
	query := `
            CREATE TABLE msg (
                id serial   NOT NULL,
                chatid int   NOT NULL,
                message char(100)   NOT NULL,
                createdat timestamp   NOT NULL,
                CONSTRAINT pk_msg PRIMARY KEY (
                    id,chatid
                 )
            );
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *chatStore) InsertChat(elchat *Chat) error {
	query := `insert into Chat 
    (personaid,createdat)Message,date,personaid
    values ($1,$2)
    `
	resp, err := s.db.Query(
		query,
		&elchat.PersonaID,
		&elchat.CreatedAt)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *chatStore) InsertMsg(elmsg *Msg) error {
	query := `insert into msg 
    (chatid,message,createdat)
    values ($1,$2,$3)
    `
	resp, err := s.db.Query(
		query,
		&elmsg.ChatID,
		&elmsg.Message,
		&elmsg.CreatedAt)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *chatStore) DeleteChat(int) error {
	return nil
}

func (s *chatStore) UpdateChat(*Chat) error {
	return nil
}

// func (s *chatStore) GetChatById(int) (*Chat, error) {
// 	return nil, nil
// }

func (s *chatStore) GetChatsByPersonaId(id int) ([]Chat, error) {
	rows, err := s.db.Query(`select * from Chats where personaid = $1`, id)

	if err == sql.ErrNoRows {

		slog.Info("no chat found with this email", "email", id)
		return nil, fmt.Errorf("chat with email [%d] not found", id)
	}
	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.PersonaID); err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (s *chatStore) GetMsgsByChatId(id int) ([]Chat, error) {
	rows, err := s.db.Query(`

    select * from msg where chatid = $1

    `, id)

	if err == sql.ErrNoRows {

		slog.Info("no chat found with this email", "email", id)
		return nil, fmt.Errorf("chat with email [%d] not found", id)
	}
	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.PersonaID); err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
