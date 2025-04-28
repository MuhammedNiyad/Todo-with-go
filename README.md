Todo API (PostgreSQL)
A RESTful API backend for a Todo application written in Go, using PostgreSQL for storage.
Features

Create tasks (POST /tasks)
List tasks (GET /tasks)
Update task status (PATCH /tasks/{id})
Delete tasks (DELETE /tasks/{id})

Prerequisites

Go 1.21 or later
PostgreSQL 15 or later
A PostgreSQL database named todo_api (and todo_api_test for tests)

Setup

Install PostgreSQL and create databases:
createdb todo_api
createdb todo_api_test


Apply the database migration:
psql -d todo_api -f migrations/001_create_tasks.sql
psql -d todo_api_test -f migrations/001_create_tasks.sql


Clone the repository:
git clone <repository-url>
cd todo-api


Install dependencies:
go mod tidy


Set environment variables (optional):
export POSTGRES_URL="postgres://your_user:your_password@localhost:5432/todo_api?sslmode=disable"
export TEST_POSTGRES_URL="postgres://your_user:your_password@localhost:5432/todo_api_test?sslmode=disable"

If not set, defaults to postgres:password@localhost:5432.

Build and run:
go build -o todo-api ./cmd/api
./todo-api



API Endpoints

POST /tasks
Body: {"description": "string"}
Response: 201 Created with task JSON


GET /tasks
Response: 200 OK with array of tasks


PATCH /tasks/{id}
Body: {"completed": true}
Response: 204 No Content


DELETE /tasks/{id}
Response: 204 No Content



Example Usage
# Create a task
curl -X POST -H "Content-Type: application/json" -d '{"description":"Buy groceries"}' http://localhost:8080/tasks

# List tasks
curl http://localhost:8080/tasks

# Complete a task
curl -X PATCH -H "Content-Type: application/json" -d '{"completed":true}' http://localhost:8080/tasks/1

# Delete a task
curl -X DELETE http://localhost:8080/tasks/1

Running Tests
go test ./internal/store

Ensure TEST_POSTGRES_URL is set or the test database is configured.
Directory Structure

cmd/api/: API server entry point
internal/handler/: HTTP handlers
internal/store/: PostgreSQL database logic
migrations/: SQL migration scripts

Troubleshooting

Connection errors: Verify PostgreSQL is running and the connection strings match your setup.
Test failures: Ensure the todo_api_test database exists and is accessible.
Table not found: Run the migration script or check NewStore for errors.

Notes

For production, use environment variables for configuration and a migration tool like golang-migrate.
Update connection strings in main.go and store_test.go if not using environment variables.

I like the Go language so I'm gonna become the best Go-lang developer ever.

