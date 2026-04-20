package db

import (
	"fmt"
	"time"
)

const (
	queryGetTask    = `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	queryDeleteTask = `DELETE FROM scheduler WHERE id = ?`
	queryUpdateDate = `UPDATE scheduler SET date = ? WHERE id = ?`
	queryUpdateTask = `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	queryTasksAll      = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	queryTasksByDate   = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
	queryTasksBySearch = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
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
	query, args := buildTasksQuery(limit, search)

	rows, err := db.Query(query, args...)
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

func buildTasksQuery(limit int, search string) (string, []any) {
	if search == "" {
		return queryTasksAll, []any{limit}
	}

	if t, err := time.Parse("02.01.2006", search); err == nil {
		return queryTasksByDate, []any{t.Format("20060102"), limit}
	}

	like := "%" + search + "%"
	return queryTasksBySearch, []any{like, like, limit}
}

func GetTask(id string) (*Task, error) {
	task := &Task{}

	row := db.QueryRow(queryGetTask, id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	res, err := db.Exec(queryUpdateTask, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
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
	_, err := db.Exec(queryDeleteTask, id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDate(date string, id string) error {
	res, err := db.Exec(queryUpdateDate, date, id)
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
