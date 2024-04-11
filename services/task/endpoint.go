package task

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/task", h.handleGetTasks).Methods(http.MethodGet)
	router.HandleFunc("/task", h.handleCreateTask).Methods(http.MethodPost)
	router.HandleFunc("/task/{id}", h.handleGetTask).Methods(http.MethodGet)
	router.HandleFunc("/task/{id}", h.handleProgressTask).Methods(http.MethodPatch)
	router.HandleFunc("/task/{id}", h.handleDeleteTask).Methods(http.MethodDelete)
	router.HandleFunc("/task/{id}", h.handleUpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/task/concurrency", h.handleConcurrencyDemo).Methods(http.MethodPost)

}
