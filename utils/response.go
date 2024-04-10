package utils

import (
	"encoding/json"
	"net/http"

	"github.com/trsnaqe/gotask/types"
)

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	errResp := types.ErrorResponse{
		Error:      err.Error(),
		StatusCode: status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}
