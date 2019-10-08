package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates the hash of the password to be stored in database
func HashPassword(password string) (string, error) {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHashWithPassword returns true if given password matches with the hash and false otherwise
func CompareHashWithPassword(hashedPassword, password string) bool {
	hash := []byte(hashedPassword)
	pass := []byte(password)
	err := bcrypt.CompareHashAndPassword(hash, pass)
	if err != nil {
		return false
	}
	return true
}
