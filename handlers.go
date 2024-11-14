package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"
)

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("getAllTasksHandler %s\n", req.URL.Path)

	allTasks := ts.store.GetAllTasks()

	var apiResponse []ResponseTask
	for _, t := range allTasks {
		apiResponse = append(apiResponse, t.ToResponseTask())
	}

	renderJSON(w, apiResponse)
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("tagHandler %s\n", req.URL.Path)

	tag := req.PathValue("tag")
	tasks := ts.store.GetTasksByTag(tag)

	var apiResponse []ResponseTask
	for _, t := range tasks {
		apiResponse = append(apiResponse, t.ToResponseTask())
	}

	renderJSON(w, apiResponse)
}

func (ts *taskServer) getTaskByIdHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("getTaskHandler %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		log.Fatal(err)
	}

	task := ts.store.GetTask(id)

	renderJSON(w, task.ToResponseTask())
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("dueHandler %s\n", req.URL.Path)

	year, errYear := strconv.Atoi(req.PathValue("year"))
	month, errMonth := strconv.Atoi(req.PathValue("month"))
	day, errDay := strconv.Atoi(req.PathValue("day"))
	if errYear != nil || errMonth != nil || errDay != nil || month < int(time.January) || month > int(time.December) {
		log.Fatal("date is invalid")
	}

	tasks := ts.store.GetTasksByDueDate(year, month, day)
	var apiResponse []ResponseTask
	for _, t := range tasks {
		apiResponse = append(apiResponse, t.ToResponseTask())
	}

	renderJSON(w, apiResponse)
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("createTaskHandler %s\n", req.URL.Path)

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Fatal(err)
	}
	if mediatype != "application/json" {
		log.Fatal(err)
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		log.Fatal(err)
	}

	tags := Tags{Tags: rt.Tags}
	id := ts.store.CreateTask(rt.Text, tags, rt.Due)
	renderJSON(w, ResponseId{Id: id})
}
