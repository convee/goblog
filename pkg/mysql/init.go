package mysql

import (
	config2 "blog/pkg/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var (
	db *sql.DB
)

func Init() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s",
		config2.MysqlConfig.Username,
		config2.MysqlConfig.Password,
		config2.MysqlConfig.Ip,
		config2.MysqlConfig.Port,
		config2.MysqlConfig.Database,
		config2.MysqlConfig.Charset,
	)
	var err error
	if db, err = sql.Open("mysql", dsn); err != nil {
		log.Fatalln("db open failed:", err)
	}
	db.SetMaxIdleConns(config2.MysqlConfig.MaxIdle)                                     //设置闲置的连接数
	db.SetMaxOpenConns(config2.MysqlConfig.MaxOpen)                                     //设置最大打开的连接数，默认值0表示不限制
	db.SetConnMaxLifetime(time.Duration(config2.MysqlConfig.MaxLifetime) * time.Second) //设置长连接的最长使用时间（从创建时开始计算），超过该时间go会自动关闭该连接
	if err := db.Ping(); err != nil {
		log.Fatalln("db ping failed:", err)
	}
}
