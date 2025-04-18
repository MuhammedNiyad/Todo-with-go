package main

import (
    "fmt"
    "net/http"
)

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

func main() {
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome to the To-Do API!")
    })
    http.ListenAndServe(":8080", nil)
}