package mysql

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// Config mysql config.
type Config struct {
	DSN             string // write data source name.
	MaxOpenConn     int    // open pool
	MaxIdleConn     int    // idle pool
	ConnMaxLifeTime int
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *sql.DB) {

	db, err := connect(c, c.DSN)
	if err != nil {
		panic(err)
	}
	return
}

func connect(c *Config, dataSourceName string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	d.SetMaxOpenConns(c.MaxOpenConn)
	d.SetMaxIdleConns(c.MaxIdleConn)
	d.SetConnMaxLifetime(time.Duration(c.ConnMaxLifeTime))
	return d, nil
}
