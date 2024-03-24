// utils.go

package main

import (
	"time"

	"github.com/google/uuid"
)

/* Utility Functions */

// currentTime returns the current time as a string in RFC 3339 format.
func CurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// GetUserByID returns a pointer to the User with the given ID.
func GetUserByID(id string) *User {
	for _, user := range cache.Users.Users {
		if user.ID == uuid.MustParse(id) {
			return &user
		}
	}
	return nil
}

func GetRoomStubByID(id string) *RoomStub {
	for _, room := range cache.Rooms.Rooms {
		if room.ID == uuid.MustParse(id) {
			return &room
		}
	}
	return nil
}
