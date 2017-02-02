package models

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)


var db *sql.DB

func InitDB() {
    var err error
    db, err = sql.Open("mysql", "getmyparts:hast1ng$@tcp(masterdb:3306)/GetMyParts?parseTime=true")
    if err != nil {
        fmt.Println("Cant connect to db",err)
        //log.Panic(err)
    }
    fmt.Println("Connection successfull!")


    if err = db.Ping(); err != nil {
        fmt.Println("Cant connect to db",err)
        //log.Panic(err)
    }
}

