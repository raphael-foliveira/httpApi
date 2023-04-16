package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/raphael-foliveira/httpApi/handlers"
)

const PORT = ":8000"

func Run() {
	serveMux := http.NewServeMux()

	serveMux.Handle("/todos/", &handlers.Todo{})

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
