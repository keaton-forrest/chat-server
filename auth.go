// auth.go

package main

import (
	"golang.org/x/crypto/bcrypt"
)

/* Authentication Functions */

// Hash a string using bcrypt
func hashString(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	return string(bytes), err
}

// Compare a password to a hash
func compareToHash(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
