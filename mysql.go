package confmysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Host     string `env:""`
	Port     int    `env:""`
	User     string `env:""`
	Password string `env:""`
	DbName   string `env:""`

	maxOpenConns int
	maxIdleConns int

	*sql.DB
}

func (my *Mysql) SetDefaults() {
	if my.Host == "" {
		my.Host = "127.0.0.1"
	}

	if my.Port == 0 {
		my.Port = 3306
	}

	if my.maxOpenConns == 0 {
		my.maxOpenConns = 10
	}

	if my.maxIdleConns == 0 {
		my.maxIdleConns = 5
	}

}

func (my *Mysql) initial() {
	if my.DB != nil {
		return
	}

	dsnTmpl := "%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True"
	dsn := fmt.Sprintf(dsnTmpl,
		my.User, my.Password,
		my.Host, my.Port,
		my.DbName,
	)

	db, err := conn(dsn)
	if err != nil {
		panic("err")
	}

	db.SetMaxOpenConns(my.maxOpenConns)
	db.SetMaxIdleConns(my.maxIdleConns)

	my.DB = db
}

func conn(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		return db, fmt.Errorf("db ping faild: %w", err)
	}

	return db, nil
}
