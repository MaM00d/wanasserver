package chat

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"

	"Server/user/persona"
)

func (s *ElChat) createChatTabel() error {
	query := `
          CREATE TABLE IF NOT EXISTS chat (id int NOT NULL, personaid int NOT NULL, userid int NOT NULL, createdat timestamp NOT NULL, CONSTRAINT pk_chat PRIMARY KEY (id
                                                                                                                                                                         , personaid
                                                                                                                                                                         , userid));
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) dropChatTabel() error {
	query := `
          DROP TABLE IF EXISTS chat;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createpersonauserfk() error {
	query := `
          ALTER TABLE chat ADD CONSTRAINT fk_chat_personaid
            FOREIGN key(personaid, userid) REFERENCES persona(
                                                        id, userid);
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) droppersonauserfk() error {
	query := `
          ALTER TABLE chat
            DROP CONSTRAINT fk_chat_personaid;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createtriggerid() error {
	query := `
          CREATE TRIGGER trig_chat_pk
            BEFORE
            INSERT ON chat
            FOR EACH ROW EXECUTE PROCEDURE fn_trig_chat_pk();
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) droptrigid() error {
	query := `
          DROP TRIGGER IF EXISTS trig_chat_pk ON chat;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) createfunctionid() error {
	query := `
          CREATE OR REPLACE FUNCTION "fn_trig_chat_pk"() RETURNS "pg_catalog"."trigger" AS $BODY$
              begin
              new.id = (select count(*) from chat where userid=new.userid and personaid=new.personaid);
              return NEW;
              end;
              $BODY$ LANGUAGE PLPGSQL VOLATILE cost 100;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) dropfunctionid() error {
	query := `
          DROP FUNCTION IF EXISTS fn_trig_chat_pk;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElChat) InsertChat(elchat *Chat) error {
	query := `insert into Chat 
          (personaid,userid,createdat)
            VALUES ($1,$2,$3)
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
	slog.Info("inserted chat to database", "personaid: ", elchat.PersonaID)

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

	err := s.db.QueryScan(
		&chats,
		`select id from Chat where userid = $1 and personaid = $2`,
		userid,
		personaid,
	)
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return nil, s.db.NotFound
	}

	return chats, nil
}
