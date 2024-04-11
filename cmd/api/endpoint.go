package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/trsnaqe/gotask/docs"
)

func registerCommonRoutes(router *mux.Router) {

	// GetLog   get-log
	// @Summary      Check logs
	// @Description  Endpoint that serves the log file
	// @Tags         API
	// @Accept       json
	// @Produce      json
	// @Success      200  {file} 	log.txt
	// @Failure      400  {object}  types.ErrorResponse
	// @Failure      500  {object}  types.ErrorResponse
	// @Router       /log [get]
	router.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app.log")
	}).Methods(http.MethodGet)

	router.HandleFunc("/log",

		// DeleteLog   delete-log
		// @Summary      Delete logs
		// @Description  Endpoint that deletes the content of log file
		// @Tags         API
		// @Accept       json
		// @Produce      json
		// @Success      200  {object}  string
		// @Failure      400  {object}  types.ErrorResponse
		// @Failure      500  {object}  types.ErrorResponse
		// @Router       /log [delete]
		func(w http.ResponseWriter, r *http.Request) {
			//delete the content of app.log
			file, err := os.OpenFile("app.log", os.O_TRUNC, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Log file has been deleted"))
		}).Methods(http.MethodDelete)

	// Health   		health
	// @Summary      Check health
	// @Description  Endpoint to check health
	// @Tags         API
	// @Accept       json
	// @Produce      json
	// @Success      200  {object} 	string
	// @Failure      400  {object}  types.ErrorResponse
	// @Failure      500  {object}  types.ErrorResponse
	// @Router       /health [get]
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Metrics   metrics
	// @Summary      Check metrics
	// @Description  Endpoint that serves the metrics
	// @Tags         API
	// @Accept       json
	// @Produce      json
	// @Success      200  {object} 	string
	// @Failure      400  {object}  types.ErrorResponse
	// @Failure      500  {object}  types.ErrorResponse
	// @Router       /metrics [get]
	router.Handle("/metrics", promhttp.Handler())

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

}
