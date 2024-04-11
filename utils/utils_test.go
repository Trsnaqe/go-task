package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/trsnaqe/gotask/types"
)

func TestParseJSON(t *testing.T) {
	// Create a request with a JSON payload
	payload := map[string]interface{}{
		"key": "value",
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(payloadBytes))

	var target map[string]interface{}

	// Test parsing JSON
	err := ParseJSON(req, &target)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the parsed JSON matches the expected payload
	if !reflect.DeepEqual(target, payload) {
		t.Errorf("Expected payload %+v, got %+v", payload, target)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	payload := map[string]interface{}{
		"key": "value",
	}

	err := WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedBody, _ := json.Marshal(payload)
	actualBody := bytes.TrimSpace(w.Body.Bytes())

	t.Logf("Expected response body: %s", expectedBody)
	t.Logf("Actual response body: %s", actualBody)

	if !bytes.Equal(actualBody, expectedBody) {
		t.Errorf("Expected response body %s, got %s", expectedBody, actualBody)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("test error")

	WriteError(w, http.StatusInternalServerError, err)

	var errResp types.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp.Error != err.Error() {
		t.Errorf("Expected error message %s, got %s", err.Error(), errResp.Error)
	}
}

func TestGetTokenFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/?token=abc123", nil)
	req.Header.Set("Authorization", "Bearer xyz456")

	token := GetTokenFromRequest(req)

	if token != "Bearer xyz456" {
		t.Errorf("Expected token %s, got %s", "Bearer xyz456", token)
	}

	reqWithoutAuth := httptest.NewRequest("GET", "/?token=abc123", nil)

	tokenWithoutAuth := GetTokenFromRequest(reqWithoutAuth)

	if tokenWithoutAuth != "abc123" {
		t.Errorf("Expected token %s, got %s", "abc123", tokenWithoutAuth)
	}
}
