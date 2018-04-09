package db

import (
	"context"
	"database/sql"
	"espressif.com/chip/factory/config"
	"log"
	"sync"
)

var db Daler
var open sync.Once

type Col int

type Doer interface {
	Save() error
	Update(cs ...Col) error
	Invisibly() error
	Delete() error
	Valid() error

	Padding(pkey string, pvalue interface{})
	AsMap(isColumnName bool, cs ...Col) map[string]interface{}
	MarshalJSON() ([]byte, error)
	UnmarshalMap(ctx context.Context, vi interface{}, cols ...Col) ([]Col, error)
}

type Daler interface {
	init() error
	Exec(q string, args ...interface{}) (int64, int64, error)
	Query(new func(cols ...string) (Doer, []interface{}, error), asterisk bool, q string, args ...interface{}) ([]Doer, error)
	Count(q string, args ...interface{}) (int, error)
	Defer() error
}

func Open(ns string) Daler {
	open.Do(func() {
		var err error
		switch config.Cfg.DriverName {
		case "sqlite3":
			db, err = NewSqlite3()
		case "mysql":
			db, err = NewMysql()
		default:
			panic("unsupport database driver!")
		}
		if err != nil {
			panic(err)
		}
	})
	return db
}

func Defer() {
	if db != nil {
		db.Defer()
	}
}

func Exec(db *sql.DB, q string, args ...interface{}) (int64, int64, error) {
	rs, err := db.Exec(q, args...)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	id, err := rs.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	return id, rows, nil
}

func Query(db *sql.DB, new func(cols ...string) (Doer, []interface{}, error), asterisk bool, q string, args ...interface{}) ([]Doer, error) {
	rows, err := db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	cols := []string{}
	if !asterisk {
		cols, err = rows.Columns()
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	dos := []Doer{}
	for rows.Next() {
		do, dest, err := new(cols...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = rows.Scan(dest...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		dos = append(dos, do)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}
	return dos, err
}

func Count(db *sql.DB, q string, args ...interface{}) (int, error) {
	rows, err := db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()
	c := 0
	for rows.Next() {
		err := rows.Scan(&c)
		if err != nil {
			log.Println(err)
		}
		return c, err
	}
	return c, nil
}
