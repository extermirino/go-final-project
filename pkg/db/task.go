package db

import (
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]Task, error) {
	var query string
	if search != "" {
		t, err := time.Parse("02.01.2006", search)

		if err == nil {
			date := t.Format("20060102")
			query = fmt.Sprintf(`SELECT * FROM scheduler WHERE date = '%s' ORDER BY date ASC LIMIT %d`, date, limit)
		} else {
			like := "%" + search + "%"
			query = fmt.Sprintf(`SELECT * FROM scheduler WHERE title LIKE '%s' OR comment LIKE '%s' ORDER BY date ASC LIMIT %d`, like, like, limit)
		}
		
	} else {
		query = fmt.Sprintf(`SELECT * FROM scheduler ORDER BY date ASC LIMIT %d`, limit)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
