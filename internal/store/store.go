package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Task represents a single todo item.
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

// Store manages the todo list using a PostgreSQL database.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store with the given database connection and ensures the tasks table exists.
func NewStore(db *sql.DB) (*Store, error) {
	// Ensure tasks table exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			description TEXT NOT NULL,
			completed BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %w", err)
	}

	return &Store{db: db}, nil
}

// AddTask adds a new task with the given description.
func (s *Store) AddTask(description string) (Task, error) {
	if description == "" {
		return Task{}, errors.New("description cannot be empty")
	}

	var task Task
	err := s.db.QueryRow(
		`INSERT INTO tasks (description, completed, created_at)
		 VALUES ($1, $2, $3)
		 RETURNING id, description, completed, created_at`,
		description, false, time.Now(),
	).Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("failed to insert task: %w", err)
	}

	return task, nil
}

// ListTasks returns all tasks.
func (s *Store) ListTasks() ([]Task, error) {
	rows, err := s.db.Query(`SELECT id, description, completed, created_at FROM tasks ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Completed, &task.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// CompleteTask marks a task as completed by ID.
func (s *Store) CompleteTask(id int, completed bool) error {
	res, err := s.db.Exec(`UPDATE tasks SET completed = $1 WHERE id = $2`, completed, id)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// DeleteTask removes a task by ID.
func (s *Store) DeleteTask(id int) error {
	res, err := s.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}