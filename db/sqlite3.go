package db

import (
	"database/sql"
	"espressif.com/chip/factory/config"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var _ Daler = (*Sqlite3)(nil)

type Sqlite3 struct {
	db *sql.DB
}

func NewSqlite3() (Daler, error) {
	s := &Sqlite3{}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Sqlite3) init() error {
	db, err := sql.Open("sqlite3", config.Cfg.DataSourceName)
	if err != nil {
		log.Println(err)
		return err
	}
	s.db = db
	return nil
}

func (s *Sqlite3) Exec(q string, args ...interface{}) (int64, int64, error) {
	return Exec(s.db, q, args...)
}

func (s *Sqlite3) Query(new func(cols ...string) (Doer, []interface{}, error), asterisk bool, q string, args ...interface{}) ([]Doer, error) {
	return Query(s.db, new, asterisk, q, args...)
}

func (s *Sqlite3) Count(q string, args ...interface{}) (int, error) {
	return Count(s.db, q, args...)
}

func (s *Sqlite3) Defer() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
