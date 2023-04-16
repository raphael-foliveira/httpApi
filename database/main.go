package database

import (
	"database/sql"
	"fmt"
	"os"
)

var Db *sql.DB

func Get() {
	var err error
	Db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	migrate()
}

func migrate() {
	fmt.Println("Updating database...")
	_, err := Db.Exec(`
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255),
		description TEXT,
		done BOOLEAN
	);

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		password VARCHAR(255)
	);
	`)
	if err != nil {
		panic(err)
	}
}
