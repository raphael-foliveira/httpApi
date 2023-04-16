package models

import (
	"github.com/raphael-foliveira/httpApi/database"
)

type ToDo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ToDoManager struct{}

func (tm *ToDoManager) Find() ([]ToDo, error) {
	rows, err := database.Db.Query("SELECT id, title, description, done FROM todos")
	if err != nil {
		return nil, err
	}
	result := []ToDo{}
	for rows.Next() {
		var curr ToDo
		err = rows.Scan(&curr.Id, &curr.Title, &curr.Description, &curr.Done)
		if err != nil {
			return nil, err
		}
		result = append(result, curr)
	}
	return result, err
}

func (tm *ToDoManager) Retrieve(id int) (ToDo, error) {
	row := database.Db.QueryRow("SELECT id, title, description, done FROM todos WHERE id = $1", id)
	var result ToDo
	return result, row.Scan(&result.Id, &result.Title, &result.Description, &result.Done)
}

func (tm *ToDoManager) Delete(id int) (ToDo, error) {
	row := database.Db.QueryRow("DELETE FROM todos WHERE id = $1 RETURNING id, title, description, done", id)
	var result ToDo
	return result, row.Scan(&result.Id, &result.Title, &result.Description, &result.Done)
}

func (tm *ToDoManager) Create(todo ToDo) (ToDo, error) {
	row := database.Db.QueryRow("INSERT INTO todos (title, description, done) VALUES ($1, $2, $3) RETURNING id, title, description, done", todo.Title, todo.Description, false)
	var result ToDo
	return result, row.Scan(&result.Id, &result.Title, &result.Description, &result.Done)
}

func (tm *ToDoManager) Update(todo ToDo) (ToDo, error) {
	var result ToDo
	previousData, err := tm.Retrieve(todo.Id)
	if err != nil {
		return result, err
	}
	if todo.Title == "" {
		result.Title = previousData.Title
	}
	if todo.Description == "" {
		result.Description = previousData.Description
	}
	row := database.Db.QueryRow("UPDATE todos SET title = $1, description = $2, done = $3 WHERE id = $4 RETURNING id, title, description, done", todo.Title, todo.Description, todo.Done, todo.Id)
	return result, row.Scan(&result.Id, &result.Title, &result.Description, &result.Done)
}
