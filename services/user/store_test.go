package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/trsnaqe/gotask/services/auth"
	"github.com/trsnaqe/gotask/types"
)

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	email := "test@example.com"
	expectedUser := &types.User{
		ID:           1,
		Email:        email,
		Password:     "password",
		RefreshToken: nil,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "refresh_token", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password, expectedUser.RefreshToken, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE email = ?").WithArgs(email).WillReturnRows(rows)

	user, err := store.GetUserByEmail(email)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("got %v, want %v", user, expectedUser)
	}
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	userID := 1
	expectedUser := &types.User{
		ID:           userID,
		Email:        "test@example.com",
		Password:     "password",
		RefreshToken: nil,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "refresh_token", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password, expectedUser.RefreshToken, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM users WHERE id = ?").WithArgs(userID).WillReturnRows(rows)

	user, err := store.GetUserByID(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("got %v, want %v", user, expectedUser)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	email := "test@example.com"
	password := "password"

	mock.ExpectExec("INSERT INTO users \\(email, password\\) VALUES \\(\\?, \\?\\)").
		WithArgs(email, password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.CreateUser(types.User{Email: email, Password: password})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	userID := 1
	email := "new_email"

	mock.ExpectExec("UPDATE users SET email = \\?, updated_at = \\? WHERE id = ?").
		WithArgs(email, sqlmock.AnyArg(), userID). // Use sqlmock.AnyArg() for time argument
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.UpdateUser(userID, types.User{Email: email})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

}

func TestChangePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	userID := 1
	oldPassword := "old_password"
	newPassword := "new_password"
	oldHashedPassword, err := auth.HashValue(oldPassword)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "refresh_token", "created_at", "updated_at"}).
		AddRow(userID, "test@example.com", oldHashedPassword, nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT \\* FROM users WHERE id = ?").WithArgs(userID).WillReturnRows(rows)
	mock.ExpectExec("UPDATE users SET password = \\?, updated_at = \\?  WHERE id = ?").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1)).
		WillReturnError(nil)

	err = store.ChangePassword(userID, oldPassword, newPassword)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewStore(db)
	userID := 1

	mock.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.DeleteUser(userID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
}
