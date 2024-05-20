package msg

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func (s *ElMsg) createMsgTabel() error {
	query := `
            CREATE TABLE if not exists msg (
                id int unique   NOT NULL,
                useriD int   NOT NULL,
                chatid int   NOT NULL,
                personaid int   NOT NULL,
                message char(100)   NOT NULL,
                createdat timestamp   NOT NULL,
                CONSTRAINT pk_msg PRIMARY KEY (
                    id,chatid
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
            new.id = (select count(*)+1 from msg where userid=new.userid and personaid=new.personaid and chatid=new.chatid);
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
    (chatid,personaid,userid,message,createdat)
    values ($1,$2,$3,$4,$5)
    `
	err := s.db.Query(
		query,
		&elmsg.ChatID,
		&elmsg.PersonaID,
		&elmsg.UserID,
		&elmsg.Message,
		&elmsg.CreatedAt)
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

func (s *ElMsg) GetMsgs(userid, personaid, chatid int) ([]MsgView, error) {
	var msgs []MsgView
	rows := s.db.QueryScan(msgs, `select message createdat from Msgs where chatid = $1 and personaid = $2 and userid = $3 `, chatid, personaid, userid)

	if rows == fmt.Errorf("not found") {

		slog.Error("GetMsgsByUserId", "id", userid)
		return nil, fmt.Errorf("msg with id [%d] not found", userid)
	}
	return msgs, nil
}
