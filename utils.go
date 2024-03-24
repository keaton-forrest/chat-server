// utils.go

package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

func GetAdminAccount() gin.Accounts {
	// Load the .env file in the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve credentials
	adminUser := os.Getenv("ADMIN_USER")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Define authorized credentials using the loaded environment variables
	adminAccount := gin.Accounts{
		adminUser: adminPassword,
	}

	return adminAccount
}

func GetChannel(roomID, userID string) *SSEStream {
	for _, channel := range cache.Channels {
		if channel.RoomID == uuid.MustParse(roomID) && channel.UserID == uuid.MustParse(userID) {
			return &channel
		}
	}
	return nil
}

func AddChannel(roomID, userID string) {
	cache.Channels = append(cache.Channels, SSEStream{
		Stream: make(chan string),
		RoomID: uuid.MustParse(roomID),
		UserID: uuid.MustParse(userID),
	})
}

func RemoveChannel(roomID, userID string) {
	for i, channel := range cache.Channels {
		if channel.RoomID == uuid.MustParse(roomID) && channel.UserID == uuid.MustParse(userID) {
			cache.Channels = append(cache.Channels[:i], cache.Channels[i+1:]...)
			return
		}
	}
}

func GetOrCreateChannel(roomID, userID string) *SSEStream {
	channel := GetChannel(roomID, userID)
	if channel == nil {
		AddChannel(roomID, userID)
		channel = GetChannel(roomID, userID)
	}
	return channel
}

func GetChannels(roomID string) []*SSEStream {
	var channels []*SSEStream
	for _, channel := range cache.Channels {
		if channel.RoomID == uuid.MustParse(roomID) {
			channels = append(channels, &channel)
		}
	}
	return channels
}
