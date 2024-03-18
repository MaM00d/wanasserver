package chat

import (
	"Server/user/persona"
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func (s *ElChat) createChatTabel() error {
	query := `
            CREATE TABLE if not exists chat (
                id int  unique   NOT NULL,
                personaid int   NOT NULL,
                useriD int   NOT NULL,
                createdat timestamp   NOT NULL,
                CONSTRAINT pk_chat PRIMARY KEY (
                    id,personaid,useriD
                 )
            );
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) dropChatTabel() error {
	query := `
    drop table if exists Chat;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createpersonafk() error {
	query := `
    ALTER TABLE chat ADD CONSTRAINT fk_chat_personaid FOREIGN KEY(personaid)
    REFERENCES persona(id);
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) droppersonafk() error {
	query := `
    ALTER TABLE chat drop CONSTRAINT fk_chat_personaid;
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

func (s *ElChat) dropuserfk() error {
	query := `
    ALTER TABLE chat drop CONSTRAINT fk_chat_userid;
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

func (s *ElChat) droptrigid() error {
	query := `
    drop trigger if exists trig_chat_pk on chat;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createfunctionid() error {
	query := `
    CREATE OR REPLACE FUNCTION "fn_trig_chat_pk"()
    RETURNS "pg_catalog"."trigger" AS $BODY$ 
    begin
    new.id = (select count(*)+1 from chat where userid=new.userid and personaid=new.personaid);
    return NEW;
    end;
    $BODY$
    LANGUAGE plpgsql VOLATILE
    COST 100;
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

func (s *ElChat) InsertChat(elchat *Chat) error {
	query := `insert into Chat 
    (personaid,userid,createdat)
    values ($1,$2,$3)
    `
	err := s.db.Query(
		query,
		&elchat.PersonaID,
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
		&elchat.PersonaID,
		&elchat.UserID,
		&elchat.CreatedAt)

	return elchat, err
}

func (s *ElChat) GetChatsByUserId(personaid, userid int) ([]*ChatView, error) {
	var chats []*ChatView
	var pers []*persona.PersonaView

	if err := s.db.QueryScan(&pers, `select id, name from persona where userid = $1 and id = $2`, userid, personaid); err != nil {
		return nil, err
	}
	if len(pers) == 0 {
		return nil, s.PersonaDoesnotExsist
	}

	err := s.db.QueryScan(&chats, `select id from Chat where userid = $1 and personaid = $2`, userid, personaid)
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return nil, s.db.NotFound
	}

	return chats, nil
}
