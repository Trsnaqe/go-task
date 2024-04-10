package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/swag/example/basic/docs"
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

// @title           GOLANG API
// @version         1.0
// @description     This is a simple golang backend API prepared for a task.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func (s *APIServer) Run() error {
	limiter := rate.NewLimiter(5, 10)

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares.ExceptionMiddleware)
	router.Use(middlewares.LoggerMiddleware)
	router.Use(middlewares.RateLimitMiddleware(limiter))

	//redirect unrecognized routes to /swagger

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userRepository := user.NewStore(s.db)
	userService := user.NewHandler(userRepository)
	userService.RegisterRoutes(subrouter)

	taskRepository := task.NewStore(s.db)
	taskService := task.NewHandler(taskRepository)
	taskService.RegisterRoutes(subrouter)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "GOLANG API"
	docs.SwaggerInfo.Description = "This is a simple golang backend API prepared for a tassk."
	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = s.address
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	registerCommonRoutes(subrouter)

	log.Println("Server is running on", s.address)

	return http.ListenAndServe(s.address, router)
}
