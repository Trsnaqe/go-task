package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashValue("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to be not empty")
	}

	if hash == "password" {
		t.Error("expected hash to be different from password")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashValue("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !CompareValue(hash, "password") {
		t.Errorf("expected value to match hash")
	}
	if CompareValue(hash, "notpassword") {
		t.Errorf("expected value to not match hash")
	}
}
