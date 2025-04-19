package store

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Use environment variable for connection string, fallback to default
	connStr := os.Getenv("TEST_POSTGRES_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/todo_api_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Start a transaction for test isolation
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		t.Fatalf("Failed to start transaction: %v", err)
	}

	// Create tasks table within transaction
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			description TEXT NOT NULL,
			completed BOOLEAN NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		tx.Rollback()
		db.Close()
		t.Fatalf("Failed to create tasks table: %v", err)
	}

	// Return db and cleanup function
	return db, func() {
		tx.Rollback()
		db.Close()
	}
}

func TestAddTask(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	task, err := store.AddTask("Test task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if task.Description != "Test task" {
		t.Errorf("Expected description 'Test task', got '%s'", task.Description)
	}
	if task.ID == 0 {
		t.Error("Expected non-zero task ID")
	}
}

func TestListTasks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	_, err = store.AddTask("Test task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestCompleteTask(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	task, err := store.AddTask("Test task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	err = store.CompleteTask(task.ID, true)
	if err != nil {
		t.Fatalf("CompleteTask failed: %v", err)
	}

	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if !tasks[0].Completed {
		t.Error("Task should be completed")
	}
}

func TestDeleteTask(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	task, err := store.AddTask("Test task")
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	err = store.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	tasks, err := store.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}