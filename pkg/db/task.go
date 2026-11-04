package db

import (
	"database/sql"
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
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get task id: %w", err)
	}

	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	query := `SELECT * FROM scheduler`
	args := make([]any, 0, 3)

	if date, err := time.Parse("02.01.2006", search); err == nil {
		query += ` WHERE date = ?`
		args = append(args, date.Format("20060102"))
	} else {
		query += ` WHERE title LIKE ? OR comment LIKE ?`
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern)
	}
	query += ` ORDER BY date LIMIT ?`
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := new(Task)
		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
