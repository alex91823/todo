package main

import (
	"encoding/json"
	"log"
	"time"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Tags string `json:"tags"`
	Due  int64  `json:"due"`
}

func (t Task) ToResponseTask() ResponseTask {

	var tags Tags
	err := json.Unmarshal([]byte(t.Tags), &tags)
	if err != nil {
		log.Fatal(err)
	}

	return ResponseTask{
		ID:   t.ID,
		Text: t.Text,
		Tags: tags.Tags,
		Due:  time.Unix(t.Due, 0).Format("2006-01-02"),
	}
}

type RequestTask struct {
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type ResponseId struct {
	Id int64 `json:"id"`
}

type Tags struct {
	Tags []string `json:"tags"`
}

type ResponseTask struct {
	ID   int      `json:"id"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
	Due  string   `json:"due"`
}
