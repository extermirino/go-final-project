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

func GetTask(id string) (*Task, error) {
	task := &Task{}

	query := fmt.Sprintf(`SELECT * FROM scheduler WHERE id = %s`, id)
	row := db.QueryRow(query)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	query := fmt.Sprintf(`UPDATE scheduler SET date = '%s', title = '%s', comment = '%s', repeat = '%s' WHERE id = %s`, task.Date, task.Title, task.Comment, task.Repeat, task.ID)

	res, err := db.Exec(query)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

func DeleteTask(id string) error {
	query := fmt.Sprintf(`DELETE FROM scheduler WHERE id = %s`, id)

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDate(date string, id string) error {
	query := fmt.Sprintf(`UPDATE scheduler SET date = '%s' WHERE id = %s`, date, id)

	res, err := db.Exec(query)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}
