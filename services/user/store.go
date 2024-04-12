package user

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/trsnaqe/gotask/services/auth"
	"github.com/trsnaqe/gotask/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	u := new(types.User)
	err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.RefreshToken, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) CreateUser(u types.User) error {
	_, err := s.db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", u.Email, u.Password)

	return err
}

func (s *Store) UpdateUser(userID int, updates types.UpdateUserPayload) error {
	var setValues []string
	var args []interface{}

	if updates.Email != nil {
		setValues = append(setValues, "email = ?")
		args = append(args, updates.Email)
	}
	if updates.Password != nil {
		setValues = append(setValues, "password = ?")
		args = append(args, updates.Password)
	}
	if updates.RefreshToken != nil {
		setValues = append(setValues, "refresh_token = ?")
		args = append(args, updates.RefreshToken)
	}
	setValues = append(setValues, "updated_at = ?")
	args = append(args, time.Now()) // current timestamp

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setValues, ", "))
	args = append(args, userID)

	_, err := s.db.Exec(query, args...)
	return err
}

func (s *Store) ChangePassword(userID int, oldPassword string, newPassword string) error {
	u, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}
	compare := auth.CompareValue(u.Password, oldPassword)
	if !compare {
		return fmt.Errorf("old password or email is incorrect")
	}
	hashedPassword, err := auth.HashValue(newPassword)
	if err != nil {
		return err
	}
	err = s.UpdateUser(u.ID, types.UpdateUserPayload{Password: &hashedPassword})
	return err
}

func (s *Store) DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := s.db.Exec(query, userID)
	return err
}
