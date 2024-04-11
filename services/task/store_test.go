package task

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/trsnaqe/gotask/types"
)

func taskComparator(task1, task2 *types.Task) bool {
	return task1.ID == task2.ID &&
		task1.Title == task2.Title &&
		task1.Description == task2.Description &&
		task1.Status == task2.Status
}

func TestGetTaskByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	expectedTask := &types.Task{
		ID:          1,
		Title:       "Task 1",
		Description: "Description for task 1",
		Status:      types.StatusPending,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	// Mock the database query
	mock.ExpectQuery("SELECT \\* FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Description, expectedTask.Status, expectedTask.CreatedAt, expectedTask.UpdatedAt))

	// Call the GetTaskByID function
	resultTask, err := store.GetTaskByID(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Check if the returned task matches the expected task
	if !reflect.DeepEqual(resultTask, expectedTask) {
		t.Errorf("expected task %+v, got %+v", expectedTask, resultTask)
	}
}

func TestGetTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	expectedTasks := []*types.Task{
		{
			ID:          1,
			Title:       "Task 1",
			Description: "Description for task 1",
			Status:      types.StatusPending,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          2,
			Title:       "Task 2",
			Description: "Description for task 2",
			Status:      types.StatusInProgress,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	mock.ExpectQuery("SELECT \\* FROM tasks").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(expectedTasks[0].ID, expectedTasks[0].Title, expectedTasks[0].Description, expectedTasks[0].Status, expectedTasks[0].CreatedAt, expectedTasks[0].UpdatedAt).
			AddRow(expectedTasks[1].ID, expectedTasks[1].Title, expectedTasks[1].Description, expectedTasks[1].Status, expectedTasks[1].CreatedAt, expectedTasks[1].UpdatedAt))

	resultTasks, err := store.GetTasks()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(resultTasks, expectedTasks) {
		for i, expected := range expectedTasks {
			if !taskComparator(&resultTasks[i], expected) {
				t.Errorf("expected task %+v, got %+v", expected, resultTasks[i])
			}
		}
	}

}

func TestGetTasksByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	status := types.StatusInProgress
	expectedTasks := []*types.Task{
		{
			ID:          1,
			Title:       "Task 1",
			Description: "Description for task 1",
			Status:      types.StatusInProgress,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			ID:          2,
			Title:       "Task 2",
			Description: "Description for task 2",
			Status:      types.StatusInProgress,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	mock.ExpectQuery("SELECT \\* FROM tasks WHERE status = ?").
		WithArgs(status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(expectedTasks[0].ID, expectedTasks[0].Title, expectedTasks[0].Description, expectedTasks[0].Status, expectedTasks[0].CreatedAt, expectedTasks[0].UpdatedAt).
			AddRow(expectedTasks[1].ID, expectedTasks[1].Title, expectedTasks[1].Description, expectedTasks[1].Status, expectedTasks[1].CreatedAt, expectedTasks[1].UpdatedAt))

	resultTasks, err := store.GetTasksByStatus(status)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(resultTasks, expectedTasks) {
		for i, expected := range expectedTasks {
			if !taskComparator(&resultTasks[i], expected) {
				t.Errorf("expected task %+v, got %+v", expected, resultTasks[i])
			}
		}
	}
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	newTask := types.Task{
		Title:       "New Task",
		Description: "Description for new task",
		Status:      types.StatusPending,
	}

	mock.ExpectExec("INSERT INTO tasks \\(title, description, status\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(newTask.Title, newTask.Description, newTask.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.CreateTask(newTask)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	taskID := 1
	updates := types.Task{
		Title: "Updated Title",
	}

	mock.ExpectExec("UPDATE tasks SET title = \\?, updated_at = \\? WHERE id = ?").
		WithArgs(updates.Title, sqlmock.AnyArg(), taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateTask(taskID, updates)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}

func TestProgressTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	taskID := 999
	expectedTask := &types.Task{
		ID:          taskID,
		Title:       "Task 1",
		Description: "Description for task 1",
		Status:      types.StatusPending,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	mock.ExpectQuery("SELECT \\* FROM tasks WHERE id = ?").
		WithArgs(taskID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"}).
			AddRow(expectedTask.ID, expectedTask.Title, expectedTask.Description, expectedTask.Status, expectedTask.CreatedAt, expectedTask.UpdatedAt))

	//it should insert one step further as the status is updated
	expectedTask.Status = types.StatusInProgress

	mock.ExpectExec("UPDATE tasks SET status = \\?, updated_at = \\? WHERE id = ?").
		WithArgs(expectedTask.Status, sqlmock.AnyArg(), taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.ProgressTask(taskID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	taskID := 1

	mock.ExpectExec("DELETE FROM tasks WHERE id = ?").
		WithArgs(taskID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.DeleteTask(taskID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}
