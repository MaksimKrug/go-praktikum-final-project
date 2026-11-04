package db

import (
	"database/sql"
	"errors"
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
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
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

func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("Task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
			  SET date = :date,
			  	  title = :title,
				  comment = :comment,
				  repeat = :repeat
			  WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("No task with such id: %s", id)
	}
	return nil
}

func UpdateDate(next, id string) error {
	query := `UPDATE scheduler
			  SET date = :date
			  WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id),
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}
