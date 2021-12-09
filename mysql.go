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
	Extra    string `env:""`

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

	if my.Extra == "" {
		my.Extra = "charset=utf8mb4&parseTime=True"
	}
}

func (my *Mysql) Init() {
	my.initial()
}

func (my *Mysql) initial() {
	if my.DB != nil {
		return
	}

	my.SetDefaults()

	db, err := my.conn()
	if err != nil {
		err = fmt.Errorf("dsn: %s, err: %v", my.dsn(), err)
		panic(err)
	}

	my.DB = db
}

func (my *Mysql) conn() (*sql.DB, error) {

	dsn := my.dsn()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return db, fmt.Errorf("db ping faild: %w", err)
	}

	db.SetMaxOpenConns(my.maxOpenConns)
	db.SetMaxIdleConns(my.maxIdleConns)

	return db, nil
}

func (my *Mysql) dsn() string {
	dsnTmpl := "%s:%s@tcp(%s:%d)/%s?%s"
	dsn := fmt.Sprintf(dsnTmpl,
		my.User, my.Password,
		my.Host, my.Port, my.DbName,
		my.Extra,
	)

	return dsn
}
