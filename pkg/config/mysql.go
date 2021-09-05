package config

import "github.com/spf13/viper"

type mysqlConfig struct {
	Ip          string
	Port        string
	Username    string
	Password    string
	Database    string
	Charset     string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime int
}

func loadDbConf() {
	MysqlConfig.Ip = viper.GetString("mysql.ip")
	MysqlConfig.Port = viper.GetString("mysql.port")
	MysqlConfig.Username = viper.GetString("mysql.username")
	MysqlConfig.Password = viper.GetString("mysql.password")
	MysqlConfig.Database = viper.GetString("mysql.database")
	MysqlConfig.Charset = viper.GetString("mysql.charset")
	MysqlConfig.MaxIdle = viper.GetInt("mysql.max_idle")
	MysqlConfig.MaxOpen = viper.GetInt("mysql.max_open")
	MysqlConfig.MaxLifetime = viper.GetInt("mysql.max_lifetime")
}
