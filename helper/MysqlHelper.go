package helper

import (
	"database/sql"
	"os"
)

type Model struct {
	DataSource string
	Db         *sql.DB
	Tx         *sql.Tx
}

func RegisterMysql() Model {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	return Model{DataSource: os.Getenv("MYSQL"), Db: db}
}

func (m *Model) CloseMysql() {
	err := m.Db.Close()
	if err != nil {
		panic(err)
	}
}
