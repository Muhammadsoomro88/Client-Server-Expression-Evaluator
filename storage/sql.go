package storage

import (
	"database/sql"
	"log"

	"github.com/cakemarketing/go-common/v5/settings"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connectionsql() *sql.DB {
	settings.AddConfigPath("config")
	settings.SetConfigName("local")
	err := settings.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err = sql.Open(settings.GetString("DATABASE"), settings.GetString("DATABASE_CONFIG")+"/"+settings.GetString("DATABASE_NAME"))
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func Sqlquery(newName string, x string, newRes float64, tstamp string) *sql.Rows {
	insert, err := db.Query("INSERT INTO ch2(name, expression, result, timestamp) values(?, ?, ?, ?)", newName, x, newRes, tstamp)
	if err != nil {
		log.Fatal(err)
	}
	return insert
}
