package mysql

import (
	"database/sql"

	"github.com/convee/goblog/pkg/storage/mysql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

// Init 初始化数据库
func Init(cfg *mysql.Config) *sql.DB {
	db = mysql.NewMySQL(cfg)
	return db
}

// GetDB 返回默认的数据库
func GetDB() *sql.DB {
	return db
}
