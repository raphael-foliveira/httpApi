package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/raphael-foliveira/httpApi/handlers"
)

const PORT = ":8000"
const dsn = "postgresql://postgres:123@localhost/gotodo?sslmode=disable"

func main() {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	handler := handlers.Handler{
		Db: db,
	}

	serveMux := http.NewServeMux()

	serveMux.Handle("/todos/", &handler)

	server := http.Server{
		Addr:         ":8000",
		Handler:      serveMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("About to listen on port", PORT)
	server.ListenAndServe()
}
