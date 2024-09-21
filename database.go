package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID        int64
	Date      time.Time
	Content   string
	Completed bool
}

func (a *App) initDB() error {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			content TEXT,
			completed BOOLEAN
		)
	`)
	return err
}

func (a *App) AddTodo(date time.Time, content string) error {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO todos (date, content, completed) VALUES (?, ?, ?)",
		date.Format("2006-01-02"), content, false)
	return err
}

func (a *App) GetTodos(date time.Time) ([]Todo, error) {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, content, completed FROM todos WHERE date = ?",
		date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var dateStr string
		err := rows.Scan(&todo.ID, &dateStr, &todo.Content, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todo.Date, _ = time.Parse("2006-01-02", dateStr)
		todos = append(todos, todo)
	}
	return todos, nil
}

func (a *App) UpdateTodo(id int64, completed bool) error {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE todos SET completed = ? WHERE id = ?", completed, id)
	return err
}
