package internal

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "ecoguard.db")
    if err != nil { return err }

    // Users table
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE,
        password TEXT,
        uploads TEXT
    )`)
    if err != nil { return err }

    // Alerts table
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS alerts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        region TEXT,
        event TEXT,
        description TEXT,
        date TEXT
    )`)
    if err != nil { return err }

    // Calendar table
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS calendar (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        crop TEXT,
        week INT,
        task TEXT
    )`)
    return err
}