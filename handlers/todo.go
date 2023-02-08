package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/raphael-foliveira/httpApi/models"
)

type Handler struct {
	Db *sql.DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todoIdParam := strings.TrimPrefix(r.URL.Path, "/todos")
	todoIdParam = strings.TrimPrefix(todoIdParam, "/")

	if todoIdParam != "" {
		todoId, err := strconv.Atoi(todoIdParam)
		if err != nil {
			http.Error(w, "todo id must be a number", http.StatusBadRequest)
		}
		if r.Method == http.MethodGet {
			foundTodo := h.RetrieveTodo(todoId)
			jsonTodo, _ := json.Marshal(foundTodo)
			w.Write(jsonTodo)
			return

		}

		if r.Method == http.MethodDelete {
			foundTodo := h.DeleteTodo(todoId)
			jsonTodo, _ := json.Marshal(foundTodo)
			w.Write(jsonTodo)
			return
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
			foundTodo, err := h.UpdateTodo(bodyTodo, todoId)
			updatedTodoJson, err := json.Marshal(foundTodo)
			w.Write(updatedTodoJson)
			return
		}
		return
	}

	if r.Method == http.MethodGet {
		todos := h.GetTodos()
		jsonTodos, err := json.Marshal(todos)
		if err != nil {
			panic(err)
		}
		w.Write(jsonTodos)
		return
	}

	if r.Method == http.MethodPost {
		defer r.Body.Close()

		bytes, err := io.ReadAll(r.Body)
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
		h.CreateTodo(newTodo)
		jsonResp, _ := json.Marshal(newTodo)
		w.Write(jsonResp)
		return
	}
}

func (h *Handler) RetrieveTodo(todoId int) models.Todo {
	rows, err := h.Db.Query("SELECT id, title, description, done FROM todos WHERE id = $1", todoId)
	if err != nil {
		return models.Todo{}
	}
	rows.Next()
	var (
		id                 int
		title, description string
		done               bool
	)
	err = rows.Scan(&id, &title, &description, &done)
	if err != nil {
		fmt.Println(err)
		return models.Todo{}
	}
	return models.Todo{
		Id:          id,
		Title:       title,
		Description: description,
		Done:        done,
	}
}

func (h *Handler) DeleteTodo(todoId int) models.Todo {
	rows, err := h.Db.Query("DELETE FROM todos WHERE id = $1 RETURNING *", todoId)
	if err != nil {
		fmt.Println(err)
		return models.Todo{}
	}
	rows.Next()
	var (
		id                 int
		title, description string
		done               bool
	)
	err = rows.Scan(&id, &title, &description, &done)
	if err != nil {
		fmt.Println(err)
		return models.Todo{}
	}
	return models.Todo{
		Id:          id,
		Title:       title,
		Description: description,
		Done:        done,
	}

}

func (h *Handler) GetTodos() []models.Todo {
	todos := []models.Todo{}
	rows, err := h.Db.Query("SELECT id, title, description, done FROM todos")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for rows.Next() {
		var (
			id          int
			title       string
			description string
			done        bool
		)
		err = rows.Scan(&id, &title, &description, &done)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		todos = append(todos, models.Todo{
			Id:          id,
			Title:       title,
			Description: description,
			Done:        done,
		})
	}
	return todos
}

func (h *Handler) CreateTodo(todoData models.Todo) (models.Todo, error) {
	todoData.Done = false
	rows, err := h.Db.Query("INSERT INTO todos (title, description, done) VALUES ($1, $2, $3) RETURNING id, title, description, done", todoData.Title, todoData.Description, todoData.Done)
	if err != nil {
		fmt.Println(err)
		return models.Todo{}, err
	}
	var (
		id                 int
		title, description string
		done               bool
	)
	rows.Scan(&id, &title, &description, &done)
	return models.Todo{
		Id:          id,
		Description: description,
		Title:       title,
		Done:        done,
	}, nil
}

func (h *Handler) UpdateTodo(todo models.Todo, todoId int) (models.Todo, error) {
	if todoId != todo.Id {
		return models.Todo{}, errors.New("cannot update todo because the object in body is different than the Id provided in the url")
	}
	rows, err := h.Db.Query("UPDATE todos SET title=$1, description=$2, done=$3 WHERE id=$4 RETURNING id, title, description, done", todo.Title, todo.Description, todo.Done, todoId)
	if err != nil {
		fmt.Println(err)
		return todo, err
	}
	rows.Next()
	var (
		id                 int
		title, description string
		done               bool
	)
	rows.Scan(&id, &title, &description, &done)
	return models.Todo{
		Id:          id,
		Description: description,
		Title:       title,
		Done:        done,
	}, nil
}
