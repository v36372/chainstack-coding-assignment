package infra

import (
	"chainstack/config"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var PostgreSql *gorm.DB

func InitPostgreSQL() {
	conf := config.Get().GetPostgreSQL()
	source := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Db,
	)
	var err error
	PostgreSql, err = gorm.Open("postgres", source)
	if err != nil {
		panic(err)
	}
	PostgreSql.DB().SetMaxIdleConns(conf.MaxIdleConns)
	PostgreSql.DB().SetMaxOpenConns(conf.MaxOpenConns)
	err = PostgreSql.DB().Ping()
	if err != nil {
		panic(err)
	}

	if conf.Debug {
		PostgreSql.LogMode(true)
	} else {
		PostgreSql.LogMode(false)
	}
}

func ClosePostgreSql() {
	if PostgreSql != nil {
		if err := PostgreSql.Close(); err != nil {
			fmt.Println("[ERROR] Cannot close Postgresql connection, err:", err)
		}
	}
}
