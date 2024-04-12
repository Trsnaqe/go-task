package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/trsnaqe/gotask/types"
)

func TestUser(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			Email:    "test",
			Password: "test",
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal("Error creating request")
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister).Methods("POST")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

}

// try login with invalid email
func TestUserLogin(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.LoginUserPayload{
			Email:    "test",
			Password: "test",
		}
		payloadJSON, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payloadJSON))
		if err != nil {
			t.Fatal("Error creating request")
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/login", handler.handleLogin).Methods("POST")
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

	})

}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {

	return nil, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {

	return nil, nil
}
func (m *mockUserStore) CreateUser(types.User) error {

	return nil
}

func (m *mockUserStore) UpdateUser(userID int, updates types.UpdateUserPayload) error {
	return nil
}

func (m *mockUserStore) ChangePassword(userID int, oldPassword string, newPassword string) error {
	return nil
}
