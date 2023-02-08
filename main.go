package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/raphael-foliveira/httpApi/controllers"
	"github.com/raphael-foliveira/httpApi/models"
)

type Message struct {
	Message string
}

const PORT = ":8000"
const dsn = "postgresql://postgres:123@localhost/gotodo?sslmode=disable"

func main() {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			todos := controllers.GetTodos(db)
			jsonTodos, err := json.Marshal(todos)
			if err != nil {
				panic(err)
			}
			w.Write(jsonTodos)
			return

		}

		if r.Method == http.MethodPost {
			defer r.Body.Close()

			bytes, _ := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("body:", string(bytes))
			newTodo := models.Todo{}

			err = json.Unmarshal(bytes, &newTodo)
			if err != nil {
				fmt.Println("bad request body")
				w.WriteHeader(http.StatusBadRequest)
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			controllers.CreateTodo(db, newTodo)
			jsonResp, _ := json.Marshal(newTodo)
			w.Write(jsonResp)
			return
		}

	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		splitUrl := strings.Split(r.URL.Path, "/")
		fmt.Println(splitUrl)
		if len(splitUrl) < 3 {
			http.Error(w, "you must provide a todo id", http.StatusBadRequest)
			return
		}
		todoId, err := strconv.Atoi(splitUrl[2])
		if err != nil {
			fmt.Println(err)
			http.Error(w, "todo id must be a number", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodGet {
			foundTodo := controllers.RetrieveTodo(db, todoId)
			jsonTodo, _ := json.Marshal(foundTodo)
			w.Write(jsonTodo)

		}
		if r.Method == http.MethodDelete {
			foundTodo := controllers.DeleteTodo(db, todoId)
			jsonTodo, _ := json.Marshal(foundTodo)
			w.Write(jsonTodo)
		}
		if r.Method == http.MethodPut {
			defer r.Body.Close()
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			var bodyTodo models.Todo
			json.Unmarshal(bytes, &bodyTodo)
			foundTodo, err := controllers.UpdateTodo(db, bodyTodo, todoId)
			updatedTodoJson, err := json.Marshal(foundTodo)
			w.Write(updatedTodoJson)
		}

	})

	fmt.Println("About to listen on port", PORT)
	http.ListenAndServe(PORT, nil)
}
