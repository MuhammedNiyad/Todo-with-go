package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"todo-with-go/internal/handler"
	"todo-with-go/internal/store"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Use environment variable for connection string, fallback to default
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, relying on environment variables: %v", err)
	}
	
	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/todo_api?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	store, err := store.NewStore(db)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	h := handler.NewHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("PATCH /tasks/", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}