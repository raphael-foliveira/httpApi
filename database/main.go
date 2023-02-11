package database

import "database/sql"

const dsn = "postgresql://postgres:123@localhost/gotodo?sslmode=disable"

func Start() *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db
}
