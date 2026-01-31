package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Conn struct {
	db *sql.DB
}

func Open(url string) *Conn {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	return &Conn{
		db: db,
	}
}

func (c *Conn) Close() {
	defer c.db.Close()
}
