package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to Database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL UNIQUE
	)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		fmt.Println("HATA:", err)
		panic("Could not create users table")
	}

	createEventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`
	_, err = DB.Exec(createEventTable)

	if err != nil {
		fmt.Println("HATA:", err)
		panic("Could not create events table")
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER 
		user_id INTEGER
		FOREIGN KEY(event_id) REFERENCES users(id),
		FOREIGN KEY(user_id) REFERENCES events(id),
	)
	`

	_, err = DB.Exec(createRegistrationTable)
	if err != nil {
		fmt.Println("HATA:", err)
		panic("Could not create events table")
	}
}
