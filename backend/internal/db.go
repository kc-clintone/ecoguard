package internal

import (
	"database/sql"
	"encoding/json"

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
        uploads TEXT,
        calendar TEXT
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

func InsertUser(user User) error {
    uploadsJSON, _ := json.Marshal(user.Uploads)
    calendarJSON, _ := json.Marshal(user.Calendar)
    _, err := DB.Exec("INSERT INTO users (username, password, uploads, calendar) VALUES (?, ?, ?, ?)",
        user.Username, user.Password, string(uploadsJSON), string(calendarJSON))
    return err
}

func GetUser(username string) (User, error) {
    var user User
    var uploadsJSON, calendarJSON string
    err := DB.QueryRow("SELECT username, password, uploads, calendar FROM users WHERE username = ?", username).Scan(
        &user.Username, &user.Password, &uploadsJSON, &calendarJSON)
    if err != nil { return user, err }
    json.Unmarshal([]byte(uploadsJSON), &user.Uploads)
    json.Unmarshal([]byte(calendarJSON), &user.Calendar)
    return user, nil
}

func UpdateUser(user User) error {
    uploadsJSON, _ := json.Marshal(user.Uploads)
    calendarJSON, _ := json.Marshal(user.Calendar)
    _, err := DB.Exec("UPDATE users SET password = ?, uploads = ?, calendar = ? WHERE username = ?",
        user.Password, string(uploadsJSON), string(calendarJSON), user.Username)
    return err
}

func UserExists(username string) (bool, error) {
    var count int
    err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
    return count > 0, err
}