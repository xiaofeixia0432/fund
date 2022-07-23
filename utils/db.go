package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DefaultDB *MysqlDBConf

type MysqlDBConf struct {
	Ip           string `yaml:"mysql_ip"`            // mysql服务器ip
	Port         int    `yaml:"mysql_port"`          // mysql服务器端口
	Database     string `yaml:"mysql_database"`      // 数据库名称
	User         string `yaml:"mysql_user"`          // 用户名
	Password     string `yaml:"mysql_password"`      // 密码
	Charset      string `yaml:"mysql_charset"`       // 编码格式
	MaxConntions int    `yaml:"mysql_maxconnetions"` // 最大连接数
	MaxIdles     int    `yaml:"mysql_maxidles"`      // 最大空闲连接数
}

var Db *sql.DB

// InitDB 建立mysql连接信息
func InitDB(conf *MysqlDBConf) error {
	var err error
	//str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", conf.User, conf.Password, conf.Ip, conf.Port, conf.Database, conf.Charset)
	//fmt.Println(str)
	Db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", conf.User, conf.Password, conf.Ip, conf.Port, conf.Database, conf.Charset))
	if err != nil {
		return err
	}

	Db.SetMaxIdleConns(conf.MaxIdles)
	Db.SetConnMaxLifetime(-1)
	Db.SetMaxOpenConns(conf.MaxConntions)
	err = Db.Ping()
	if err != nil {
		return err
	}
	return nil
}
