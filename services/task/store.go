package task

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/trsnaqe/gotask/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}
func (s *Store) GetTaskByID(taskID int) (*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE id = ?", taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		t, err := scanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
		return t, nil
	}

	return nil, errors.New("no task found with the given ID")
}

func (s *Store) GetTasks() ([]types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	tasks := make([]types.Task, 0)
	for rows.Next() {
		t, err := scanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, nil
}

func scanRowIntoTask(rows *sql.Rows) (*types.Task, error) {
	t := new(types.Task)
	err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Store) CreateTask(t types.Task) error {

	_, err := s.db.Exec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", t.Title, t.Description, t.Status)

	return err
}

func validateCreateTaskPayload(payload types.CreateTaskPayload) error {
	if payload.Status != types.StatusPending &&
		payload.Status != types.StatusInProgress &&
		payload.Status != types.StatusCompleted {
		return errors.New("invalid task status")
	}
	return nil
}

func (s *Store) UpdateTask(taskID int, updates types.Task) error {
	var setValues []string
	var args []interface{}

	if updates.Title != "" {
		setValues = append(setValues, "title = ?")
		args = append(args, updates.Title)
	}
	if updates.Description != "" {
		setValues = append(setValues, "description = ?")
		args = append(args, updates.Description)
	}
	if updates.Status != "" {
		setValues = append(setValues, "status = ?")
		args = append(args, updates.Status)
	}
	setValues = append(setValues, "updated_at = ?")
	args = append(args, time.Now()) // current timestamp

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(setValues, ", "))
	args = append(args, taskID)

	_, err := s.db.Exec(query, args...)
	return err
}

func (s *Store) ProgressTask(taskID int) error {
	task, err := s.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	switch task.Status {
	case types.StatusPending:
		task.Status = types.StatusInProgress
	case types.StatusInProgress:
		task.Status = types.StatusCompleted
	case types.StatusCompleted:
		return errors.New("task is already completed")
	}

	return s.UpdateTask(taskID, *task)
}

func (s *Store) DeleteTask(taskID int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	return err
}
