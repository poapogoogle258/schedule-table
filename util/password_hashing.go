package util

import (
	"os"

	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) string {
	generateFromPasswordCost := 8

	if os.Getenv("BCRYPT_COST") != "" {
		generateFromPasswordCost, _ = strconv.Atoi(os.Getenv("BCRYPT_COST"))
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), generateFromPasswordCost)
	return string(bytes)
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
