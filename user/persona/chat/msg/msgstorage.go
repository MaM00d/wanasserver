package msg

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"Server/user/persona"
)

func (s *ElMsg) createMsgTabel() error {
	query := `
            CREATE TABLE if not exists msg (
                id int NOT NULL,
                useriD int   NOT NULL,
                chatid int   NOT NULL,
                personaid int   NOT NULL,
                message char(100)   NOT NULL,
                createdat timestamp   NOT NULL,
                state boolean not null,
                CONSTRAINT pk_msg PRIMARY KEY (
                    id,personaid,chatid,userid
                 )
            );
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) dropMsgTabel() error {
	query := `
    drop table if exists Msg;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) createchatfk() error {
	query := `
    ALTER TABLE msg ADD CONSTRAINT fk_msg_chatid FOREIGN KEY(chatid,userid,personaid)
    REFERENCES chat (id,userid,personaid);
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) dropchatfk() error {
	query := `
    ALTER TABLE msg
    drop CONSTRAINT fk_msg_chatid;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) createfunctionid() error {
	query := `
            CREATE OR REPLACE FUNCTION "fn_trig_msg_pk"()
              RETURNS "pg_catalog"."trigger" AS $BODY$ 
            begin
            new.id = (select count(*) from msg where userid=new.userid and personaid=new.personaid and chatid=new.chatid);
            return NEW;
            end;
            $BODY$
              LANGUAGE plpgsql VOLATILE
              COST 100;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) dropfunctionid() error {
	query := `
    drop function if exists fn_trig_msg_pk;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) createtriggerid() error {
	query := `

            CREATE TRIGGER trig_msg_pk
              BEFORE insert 
              ON msg
              FOR EACH ROW
              EXECUTE PROCEDURE fn_trig_msg_pk();
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) droptrigid() error {
	query := `
    drop trigger if exists trig_msg_pk on msg;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElMsg) InsertMsg(elmsg *Msg) error {
	query := `insert into Msg 
    (chatid,personaid,userid,message,createdat,state)
    values ($1,$2,$3,$4,$5,$6)
    `
	err := s.db.Query(
		query,
		&elmsg.ChatID,
		&elmsg.PersonaID,
		&elmsg.UserID,
		&elmsg.Message,
		&elmsg.CreatedAt,
		&elmsg.State,
	)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	return nil
}

func (s *ElMsg) DeleteMsg(int) error {
	return nil
}

func (s *ElMsg) UpdateMsg(*Msg) error {
	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Msg, error) {
	elmsg := new(Msg)
	err := rows.Scan(
		&elmsg.ChatID,
		&elmsg.PersonaID,
		&elmsg.UserID,
		&elmsg.Message,
		&elmsg.CreatedAt)

	return elmsg, err
}

func (s *ElMsg) GetMsgsByUserId(id int) ([]Msg, error) {
	var msgs []Msg
	rows := s.db.QueryScan(msgs, `select * from Msgs where userid = $1`, id)

	if rows == fmt.Errorf("not found") {

		slog.Error("GetMsgsByUserId", "id", id)
		return nil, fmt.Errorf("msg with id [%d] not found", id)
	}
	return msgs, nil
}

func (s *ElMsg) GetMsgs(userid, personaid, chatid int) (*[]MsgView, error) {
	var msgs []MsgView
	var pers []persona.PersonaView

	if err := s.db.QueryScan(&pers, `select id, name from persona where userid = $1 and id = $2`, userid, personaid); err != nil {
		return nil, err
	}
	if len(pers) == 0 {
		return nil, s.PersonaDoesnotExsist
	}

	err := s.db.QueryScan(
		&msgs,
		`select message ,state,createdat from msg where chatid = $1 and personaid = $2 and userid = $3`,
		chatid,
		personaid,
		userid,
	)
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, s.db.NotFound
	}
	slog.Info("done", "msgs:", msgs)
	return &msgs, nil
}
