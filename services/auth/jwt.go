package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trsnaqe/gotask/config"
	"github.com/trsnaqe/gotask/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccessToken(userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTAccessExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiresAt": time.Now().Add(expiration).Unix(),
		"type":      "access",
	})

	tokenString, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTRefreshExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiresAt": time.Now().Add(expiration).Unix(),
		"type":      "refresh",
	})

	tokenString, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateTokens(userID int) (tokens types.Tokens, err error) {
	accessToken, err := CreateAccessToken(userID)
	if err != nil {
		return types.Tokens{AccessToken: "", RefreshToken: ""}, err
	}

	refreshToken, err := CreateRefreshToken(userID)
	if err != nil {
		return types.Tokens{AccessToken: "", RefreshToken: ""}, err
	}

	return types.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(types.UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}

func HashRefreshToken(value string) (string, error) {
	shaHash := sha256.New()
	shaHash.Write([]byte(value))
	hashedValue := shaHash.Sum(nil)

	bcryptHash, err := bcrypt.GenerateFromPassword(hashedValue, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcryptHash), nil
}

func CompareRefreshToken(hashedValue, value string) bool {
	shaHash := sha256.New()
	shaHash.Write([]byte(value))
	hashedInput := shaHash.Sum(nil)

	return CompareValue(hashedValue, string(hashedInput))
}
