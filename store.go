package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func CreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT, tags TEXT, due INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
}

func SeedData(db *sql.DB) {
	_, err := db.Exec("INSERT INTO tasks(text, tags, due) VALUES(?, ?, ?)", "task-1", `{"tags":["tag1", "tag2"]}`, time.Now().AddDate(0, 0, 1).Unix())
	_, err = db.Exec("INSERT INTO tasks(text, tags, due) VALUES(?, ?, ?)", "task-2", `{"tags":["tag2", "tag3"]}`, time.Now().AddDate(0, 0, 2).Unix())
	_, err = db.Exec("INSERT INTO tasks(text, tags, due) VALUES(?, ?, ?)", "task-3", `{"tags":["tag4", "tag1"]}`, time.Now().AddDate(0, 0, 3).Unix())
	_, err = db.Exec("INSERT INTO tasks(text, tags, due) VALUES(?, ?, ?)", "task-4", `{"tags":["tag3", "tag1"]}`, time.Now().AddDate(0, 0, 4).Unix())
	if err != nil {
		log.Fatal(err)
	}
}

func (ts *TaskStore) GetAllTasks() []Task {
	rows, err := ts.db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
	}

	var allTasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Text, &t.Tags, &t.Due)
		if err != nil {
			log.Fatal(err)
		}
		allTasks = append(allTasks, t)
	}
	return allTasks
}

func (ts *TaskStore) GetTasksByTag(tag string) []Task {

	// https://sqlite.org/json1.html
	// https://stackoverflow.com/questions/74731854/how-to-check-if-array-contains-an-item-in-json-column-using-sqlite
	stmt :=
		`SELECT *
		FROM tasks
		WHERE CASE
		WHEN tags LIKE '{"tags":[%]}' THEN
		EXISTS (
		SELECT *
		FROM json_each(json_extract(tasks.tags,'$.tags'))
		WHERE json_each.value = ?
		)
		END;
	`

	rows, err := ts.db.Query(stmt, tag)
	if err != nil {
		log.Fatal(err)
	}

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Text, &t.Tags, &t.Due)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}
	return tasks
}

func (ts *TaskStore) GetTask(id int) Task {
	rows, err := ts.db.Query("SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	var t Task
	for rows.Next() {
		rows.Scan(&t.ID, &t.Text, &t.Tags, &t.Due)
		if err != nil {
			log.Fatal(err)
		}
	}
	return t
}

func (ts *TaskStore) GetTasksByDueDate(year int, month int, day int) []Task {

	t, err := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", year, month, day))
	if err != nil {
		log.Fatal(err)
	}

	y, m, d := t.AddDate(0, 0, 1).Date()
	stmt := fmt.Sprintf(`SELECT * FROM tasks WHERE due < unixepoch('%d-%d-%d')`, y, m, d)
	rows, err := ts.db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Text, &t.Tags, &t.Due)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}
	return tasks
}

func (ts *TaskStore) CreateTask(text string, tags Tags, due time.Time) int64 {

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		log.Fatal(err)
	}

	// jsonTags := fmt.Sprintf(`{"tags":["%s]}`, strings.Join(tags, `,"`))
	result, err := ts.db.Exec("INSERT INTO tasks (text, tags, due) VALUES (?, ?, ?)", text, string(jsonTags), due.Unix())
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	return id
}
