package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/gotask/middlewares"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
	router.HandleFunc("/logout", middlewares.AuthMiddleware(h.handleLogout, h.store)).Methods(http.MethodPost)
	router.HandleFunc("/refresh", middlewares.AuthMiddleware(h.handleRefreshToken, h.store)).Methods(http.MethodPost)
	router.HandleFunc("/reset-password", middlewares.AuthMiddleware(h.ChangePassword, h.store)).Methods(http.MethodPost)
}
