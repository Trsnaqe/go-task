package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/trsnaqe/gotask/types"
)

func TestTask(t *testing.T) {
	taskStore := &mockTaskStore{}
	handler := NewHandler(taskStore)

	t.Run("should create a task with valid payload", func(t *testing.T) {
		payload := types.CreateTaskPayload{
			Title:       "Task 1",
			Description: "Description of Task 1",
			Status:      types.StatusPending,
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal("Error creating request")
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/task", handler.handleCreateTask).Methods("POST")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should update a task with valid payload", func(t *testing.T) {
		payload := types.Task{
			Title:       "Updated Task 1",
			Description: "Updated Description of Task 1",
			Status:      types.StatusInProgress,
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("PUT", "/task/1", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal("Error creating request")
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/task/{id}", handler.handleUpdateTask).Methods("PUT")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should delete a task with valid ID", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/task/1", nil)
		if err != nil {
			t.Fatal("Error creating request")
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/task/{id}", handler.handleDeleteTask).Methods("DELETE")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should create a task with valid payload", func(t *testing.T) {
		payload := types.CreateTaskPayload{
			Title:       "Task 3",
			Description: "Description of Task 3",
			Status:      types.StatusPending,
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(payloadJSON))
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/task", handler.handleCreateTask).Methods("POST")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should update a task with valid payload", func(t *testing.T) {
		payload := types.Task{
			Title:       "Updated Task 1",
			Description: "Updated Description of Task 1",
			Status:      types.StatusInProgress,
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("PUT", "/task/1", bytes.NewBuffer(payloadJSON))
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/task/{id}", handler.handleUpdateTask).Methods("PUT")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should delete a task with valid ID", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/task/1", nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/task/{id}", handler.handleDeleteTask).Methods("DELETE")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

}

type mockTaskStore struct{}

func (m *mockTaskStore) GetTasks() ([]types.Task, error) {
	return nil, nil
}

func (m *mockTaskStore) GetTaskByID(id int) (*types.Task, error) {
	return nil, nil
}

func (m *mockTaskStore) CreateTask(task types.Task) error {
	return nil
}

func (m *mockTaskStore) UpdateTask(taskID int, updates types.Task) error {
	return nil
}

func (m *mockTaskStore) DeleteTask(taskID int) error {
	return nil
}

func (m *mockTaskStore) ProgressTask(taskID int) error {
	return nil
}
func (m *mockTaskStore) GetTasksByStatus(status types.TaskStatus) ([]types.Task, error) {
	return nil, nil
}
