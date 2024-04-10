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

	router.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app.log")
	}).Methods(http.MethodGet)

	// @Summary      Get log file
	// @Description  Retrieve the log file
	// @Produce      plain
	// @Success      200  {string}  string
	router.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
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

	// @Summary      Show an account

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

}
