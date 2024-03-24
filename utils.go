// utils.go

package main

import (
	"time"

	"github.com/google/uuid"
)

/* Utility Functions */

// currentTime returns the current time as a string in RFC 3339 format.
func currentTime() string {
	return time.Now().Format(time.RFC3339)
}

// GetUserByID returns a pointer to the User with the given ID.
func getUserByID(id string) *User {
	for _, user := range cache.Users.Users {
		if user.ID == uuid.MustParse(id) {
			return &user
		}
	}
	return nil
}
