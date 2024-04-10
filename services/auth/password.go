package auth

import "golang.org/x/crypto/bcrypt"

func HashValue(value string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareValue(hashedValue, value string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value)) == nil
}
