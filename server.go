package main

import "database/sql"

type taskServer struct {
	store *TaskStore
}

func NewTaskServer(db *sql.DB) *taskServer {
	return &taskServer{store: NewTaskStore(db)}
}
