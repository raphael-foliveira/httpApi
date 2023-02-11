package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/raphael-foliveira/httpApi/handlers"
)

const PORT = ":8000"

func Run(db *sql.DB) {
	todoHandler := handlers.Todo{
		Db: db,
	}

	serveMux := http.NewServeMux()

	serveMux.Handle("/todos/", &todoHandler)

	server := http.Server{
		Addr:         ":8000",
		Handler:      serveMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	defer server.Shutdown(context.Background())

	fmt.Println("About to listen on port", PORT)
	server.ListenAndServe()
}
