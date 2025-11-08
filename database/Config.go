package Config

import (
    "database/sql"
    "fmt"
    _ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() error {
    var err error
    DB, err = sql.Open("sqlite", "./mydb.db")
    if err != nil {
        return err
    }
    
    // Test the connection
    if err = DB.Ping(); err != nil {
        return err
    }
    
    fmt.Println("Connected to SQLite!")
    return nil
}