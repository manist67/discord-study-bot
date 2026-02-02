package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Conn struct {
	db *sqlx.DB
}

func Open(url string) *Conn {
	db, err := sqlx.Connect("mysql", url)
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
