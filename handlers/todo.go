package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/raphael-foliveira/httpApi/models"
)

type Todo struct{}

var manager models.ToDoManager

func (h *Todo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoIdParam := strings.TrimPrefix(r.URL.Path, "/todos")
	todoIdParam = strings.TrimPrefix(todoIdParam, "/")

	if todoIdParam != "" {
		todoId, err := strconv.Atoi(todoIdParam)
		if err != nil {
			http.Error(w, "todo id must be a number", http.StatusBadRequest)
		}

		if r.Method == http.MethodGet {
			toDo, err := manager.Retrieve(todoId)
			if err != nil {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(toDo)
			return
		}

		if r.Method == http.MethodDelete {
			foundTodo, err := manager.Delete(todoId)
			if err != nil {
				http.Error(w, "could not delete", http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(foundTodo)
			return
		}

		if r.Method == http.MethodPut {
			defer r.Body.Close()
			var bodyTodo models.ToDo
			err = json.NewDecoder(r.Body).Decode(&bodyTodo)
			if err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			bodyTodo.Id = todoId
			foundTodo, err := manager.Update(bodyTodo)
			if err != nil {
				http.Error(w, "could not update", http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(foundTodo)
			return
		}
		return
	}

	if r.Method == http.MethodGet {
		todos, err := manager.Find()
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "error retrieving todos"})
			return
		}
		json.NewEncoder(w).Encode(todos)
		return
	}

	if r.Method == http.MethodPost {
		defer r.Body.Close()
		newTodo := models.ToDo{}
		err := json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil {
			http.Error(w, "request body invalidated", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		createdTodo, err := manager.Create(newTodo)
		if err != nil {
			http.Error(w, "could not create todo", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(createdTodo)
		return
	}
}
