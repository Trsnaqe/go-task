package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseJSON(r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return errors.New("please send a request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
