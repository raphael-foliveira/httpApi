package controllers

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/raphael-foliveira/httpApi/models"
)

func RetrieveTodo(db *sql.DB, todoId int) models.Todo {
	rows, err := db.Query("SELECT id, title, description, done FROM todos WHERE id = $1", todoId)
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

func DeleteTodo(db *sql.DB, todoId int) models.Todo {
	rows, err := db.Query("DELETE FROM todos WHERE id = $1 RETURNING *", todoId)
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

func GetTodos(db *sql.DB) []models.Todo {
	todos := []models.Todo{}
	rows, err := db.Query("SELECT id, title, description, done FROM todos")
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

func CreateTodo(db *sql.DB, todoData models.Todo) (models.Todo, error) {
	todoData.Done = false
	rows, err := db.Query("INSERT INTO todos (title, description, done) VALUES ($1, $2, $3) RETURNING id, title, description, done", todoData.Title, todoData.Description, todoData.Done)
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

func UpdateTodo(db *sql.DB, todo models.Todo, todoId int) (models.Todo, error) {
	if todoId != todo.Id {
		return models.Todo{}, errors.New("cannot update todo because the object in body is different than the Id provided in the url")
	}
	rows, err := db.Query("UPDATE todos SET title=$1, description=$2, done=$3 WHERE id=$4 RETURNING id, title, description, done", todo.Title, todo.Description, todo.Done, todoId)
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
