package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/gotask/middlewares"
	"github.com/trsnaqe/gotask/services/task"
	"github.com/trsnaqe/gotask/services/user"
	"golang.org/x/time/rate"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}

}

// @title                      GOLANG API
// @description                This is a simple Golang backend API prepared for a task.
// @version                    1.0
//
// @contact                    {
//	  name: "API Support",
//	  url: "http://www.example.com/support",
//	  email: "support@example.com"
//	}
//
// @license                    {
//	  name: "Apache 2.0",
//	  url: "http://www.apache.org/licenses/LICENSE-2.0.html"
//	}
//
// @host                       localhost:8080
// @BasePath                   /api/v1
// @securityDefinitions.apiKey jwtKey
// @in                         header
// @name                       Authorization
func (s *APIServer) Run() error {
	limiter := rate.NewLimiter(5, 10)

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares.ExceptionMiddleware)
	router.Use(middlewares.LoggerMiddleware)
	router.Use(middlewares.RateLimitMiddleware(limiter))

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userRepository := user.NewStore(s.db)
	userService := user.NewHandler(userRepository)
	userService.RegisterRoutes(subrouter)

	taskRepository := task.NewStore(s.db)
	taskService := task.NewHandler(taskRepository)
	taskService.RegisterRoutes(subrouter)

	registerCommonRoutes(subrouter)

	log.Println("Server is running on", s.address)

	return http.ListenAndServe(s.address, router)
}
