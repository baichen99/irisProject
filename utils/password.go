package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword returns hash value of a string
func HashPassword(password string) (hashedPassword string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword = string(hash)
	return
}

// ComparePassword returns check password result
func ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
