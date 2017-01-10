package models

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)


var db *sql.DB

func InitDB() {
    var err error
    db, err = sql.Open("mysql", "webapi:hast1ng$@tcp(fmpsql.cloudapp.net:3306)/fmp?parseTime=true")
    if err != nil {
        log.Panic(err)
    }

    if err = db.Ping(); err != nil {
        log.Panic(err)
    }
}