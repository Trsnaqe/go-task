package auth

import (
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateJWT(t *testing.T) {
	tokens, err := CreateTokens(123)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("expected access token to be not empty")
	}
	if tokens.RefreshToken == "" {
		t.Error("expected refresh token to be not empty")
	}

}

func TestValidateJWT(t *testing.T) {
	userID := 123
	tokenString, err := CreateAccessToken(userID)
	if err != nil {
		t.Fatalf("error creating access token: %v", err)
	}

	token, err := ValidateJWT(tokenString)
	if err != nil {
		t.Errorf("error validating JWT: %v", err)
	}

	if !token.Valid {
		t.Error("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims are not of type MapClaims")
	}
	claimUserID, ok := claims["userID"].(string)
	if !ok {
		t.Fatal("userID claim is not a string")
	}

	parsedUserID, err := strconv.Atoi(claimUserID)
	if err != nil {
		t.Fatalf("error parsing userID from claim: %v", err)
	}

	if parsedUserID != userID {
		t.Errorf("expected userID to be %d, got %d", userID, parsedUserID)
	}
}

func TestHashRefreshToken(t *testing.T) {
	value := "myrefresh"
	hashedValue, err := HashRefreshToken(value)
	if err != nil {
		t.Fatalf("error hashing refresh token: %v", err)
	}

	if hashedValue == value {
		t.Error("expected hashed value to be different from original value")
	}
}

func TestCompareRefreshToken(t *testing.T) {
	// Test CompareRefreshToken
	value := "myrefresh"
	hashedValue, err := HashRefreshToken(value)
	if err != nil {
		t.Fatalf("error hashing refresh token: %v", err)
	}

	if !CompareRefreshToken(hashedValue, value) {
		t.Error("expected comparison to return true for the same value")
	}

	if CompareRefreshToken(hashedValue, "anotherrefresh") {
		t.Error("expected comparison to return false for a different value")
	}
}
