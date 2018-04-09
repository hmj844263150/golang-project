package db

import (
	"database/sql"
	"espressif.com/chip/factory/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var _ Daler = (*Mysql)(nil)

type Mysql struct {
	db *sql.DB
}

func NewMysql() (Daler, error) {
	m := &Mysql{}
	err := m.init()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Mysql) init() error {
	db, err := sql.Open(config.Cfg.DriverName, config.Cfg.DataSourceName+"?charset=utf8&loc=Asia%2FShanghai&parseTime=true")
	if err != nil {
		log.Println(err)
		return err
	}
	m.db = db
	m.db.SetConnMaxLifetime(2 * time.Hour)
	return nil
}

func (m *Mysql) Exec(q string, args ...interface{}) (int64, int64, error) {
	return Exec(m.db, q, args...)
}

func (m *Mysql) Query(new func(cols ...string) (Doer, []interface{}, error), asterisk bool, q string, args ...interface{}) ([]Doer, error) {
	return Query(m.db, new, asterisk, q, args...)
}

func (m *Mysql) Count(q string, args ...interface{}) (int, error) {
	return Count(m.db, q, args...)
}

func (m *Mysql) Defer() error {
	return m.db.Close()
}
