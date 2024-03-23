package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
)

/* File Handlers */

// LoadRoom loads a room from the JSON file associated with the room ID
func LoadRoom(roomID string) (Room, error) {
	var room Room
	filename := roomID + ".json" // Construct the filename using the room ID
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return room, err
	}
	err = json.Unmarshal(data, &room)
	return room, err
}

// SaveRoom saves a room to the JSON file associated with the room ID
func SaveRoom(room Room) error {
	filename := room.ID.String() + ".json" // Use room.ID.String() for filename
	data, err := json.MarshalIndent(room, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// LoadUsers loads the users from the JSON file
func LoadUsers() (Users, error) {
	var users Users
	data, err := os.ReadFile("users.json")
	if err != nil {
		log.Println(err)
		return users, err
	}
	err = json.Unmarshal(data, &users)
	return users, err
}

func SaveUsers(users Users) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}
	return os.WriteFile("users.json", data, 0644)
}

// CreateRoomFileIfNotExist creates a new room file with the specified room ID if it doesn't exist
func CreateRoomFileIfNotExist(roomID uuid.UUID) error {
	filename := roomID.String() + ".json" // Use roomID.String() for filename

	// Write {} to the file to represent an empty json file
	err := os.WriteFile(filename, []byte("{}"), 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	// If the file already exists or has been created successfully, return nil
	return nil
}
