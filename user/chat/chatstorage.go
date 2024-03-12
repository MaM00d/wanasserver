package chat

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func (s *ElChat) InitDb() error {
	if err := s.createChatTabel(); err != nil {
		return err
	}
	if err := s.createuserfk(); err != nil {
		return err
	}

	if err := s.createfunctionid(); err != nil {
		return err
	}
	if err := s.createtriggerid(); err != nil {
		return err
	}

	return s.createtriggerid()
}

func (s *ElChat) DropDb() error {
	if err := s.droptrigid(); err != nil {
		return err
	}
	if err := s.dropfunctionid(); err != nil {
		return err
	}
	if err := s.dropChatTabel(); err != nil {
		return err
	}
	if err := s.dropuserfk(); err != nil {
		return err
	}

	return nil
}

func (s *ElChat) dropChatTabel() error {
	query := `
    drop table if exists Chat;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) dropuserfk() error {
	query := `
    drop CONSTRAINT if exists fk_chat_personaid
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createuserfk() error {
	query := `

    ALTER TABLE chat ADD CONSTRAINT fk_chat_useriD FOREIGN KEY(useriD)
    REFERENCES users (id);
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) dropfunctionid() error {
	query := `
    drop function if exists fn_trig_chat_pk;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) droptrigid() error {
	query := `
    drop trigger if exists trig_chat_pk on chat;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createChatTabel() error {
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
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createtriggerid() error {
	query := `

            CREATE TRIGGER trig_chat_pk
              BEFORE insert 
              ON chat
              FOR EACH ROW
              EXECUTE PROCEDURE fn_trig_chat_pk();
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createfunctionid() error {
	query := `
            CREATE OR REPLACE FUNCTION "fn_trig_chat_pk"()
              RETURNS "pg_catalog"."trigger" AS $BODY$ 
            begin
            new.id = (select count(*)+1 from chat where personaid=new.personaid);
            return NEW;
            end;
            $BODY$
              LANGUAGE plpgsql VOLATILE
              COST 100;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) InsertChat(elchat *Chat) error {
	query := `insert into Chat 
    (name,personaid,createdat)
    values ($1,$2,$3)
    `
	err := s.db.Query(
		query,
		&elchat.Name,
		&elchat.UserID,
		&elchat.CreatedAt)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	return nil
}

func (s *ElChat) DeleteChat(int) error {
	return nil
}

func (s *ElChat) UpdateChat(*Chat) error {
	return nil
}

// func (s *ElChat) GetChatById(int) (*Chat, error) {
// 	return nil, nil
// }

func scanIntoAccount(rows *sql.Rows) (*Chat, error) {
	elchat := new(Chat)
	err := rows.Scan(
		&elchat.ID,
		&elchat.Name,
		&elchat.UserID,
		&elchat.CreatedAt)

	return elchat, err
}

func (s *ElChat) GetChatsByUserId(id int) ([]Chat, error) {
	var chats []Chat
	rows := s.db.QueryScan(chats, `select * from Chats where ID = $1`, id)

	if rows == fmt.Errorf("not found") {

		slog.Error("GetChatsByUserId", "id", id)
		return nil, fmt.Errorf("chat with id [%d] not found", id)
	}
	return chats, nil
}
