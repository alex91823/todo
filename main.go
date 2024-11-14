package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	os.Remove("./todo.db")

	db, err := sql.Open("sqlite", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	CreateTable(db)
	SeedData(db)
	server := NewTaskServer(db)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/", server.createTaskHandler)
	mux.HandleFunc("GET /tasks/", server.getAllTasksHandler)
	mux.HandleFunc("GET /tasks/{id}/", server.getTaskByIdHandler)
	// mux.HandleFunc("DELETE /tasks/{id}/", server.deleteTaskByIdHandler)
	// mux.HandleFunc("DELETE /tasks/", server.deleteAllTasksHandler)
	mux.HandleFunc("GET /tags/{tag}/", server.tagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)

	log.Println("Starting REST API server: http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
