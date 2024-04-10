package task

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/trsnaqe/gotask/types"
	"github.com/trsnaqe/gotask/utils"
)

type Handler struct {
	store types.TaskStore
	mu    sync.Mutex
	queue chan int
}

func NewHandler(store types.TaskStore) *Handler {
	handler := &Handler{
		store: store,
		queue: make(chan int, 2), // 2 workers
	}
	handler.StartWorkers(2) // Start 2 worker goroutines
	return handler
}

// Returns tasks
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tasks)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateTaskPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err.(validator.ValidationErrors)))
		return
	}

	err = validateCreateTaskPayload(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.CreateTask(types.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Status:      payload.Status,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

// needs a valid task id
func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task ID is missing in URL"))
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	_, err = h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("Parsing JSON request body into updates struct")
	var updates types.Task
	err = utils.ParseJSON(r, &updates)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("Parsed JSON request body into updates struct")

	err = h.store.UpdateTask(taskID, updates)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Println("Task updated successfully")

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Task updated successfully"})
}

// needs a valid task id
func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task ID is missing in URL"))
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	_, err = h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.DeleteTask(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

// Progress task between stages, utilizes mutex to prevent concurrent progress thus race condition
func (h *Handler) handleProgressTask(w http.ResponseWriter, r *http.Request) {
	locked := h.mu.TryLock()
	if locked {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("another process is already in progress, please try again later"))
		return
	}
	defer h.mu.Unlock()

	vars := mux.Vars(r)
	taskIDStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("task ID is missing in URL"))
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	_, err = h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.ProgressTask(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Task %d progressed successfully", taskID)})
}

// Handle Currency demo to demonstrate queued processing
func (h *Handler) handleConcurrencyDemo(w http.ResponseWriter, r *http.Request) {
	for i := 1; i <= 4; i++ {
		// add queued task to prometheus
		// I want to see how many tasks are in the queue
		// and how many processed using total etc
		h.queue <- i // Add task to the queue
		utils.IncrementQueueLength()
		log.Printf("Task %d added to the queue", i)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Tasks added to the queue for concurrency demonstration, check logs for processing status and prometheus metrics in `api/v1/metrics` for queue length and tasks processed"})
}

func (h *Handler) worker() {
	for range h.queue {
		utils.IncrementQueueLength()
		log.Println("Processing task")
		// Simulate processing time
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
		log.Println("Task processed")
		utils.IncrementTasksProcessed()
		utils.DecrementQueueLength()
	}
}

func (h *Handler) StartWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go h.worker()
	}
}
