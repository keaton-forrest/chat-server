// cache.go

package main

import (
	"log"
)

// Cache is an interface for caching data
var cache Cache

func InitCache() {
	users, err := LoadUsers()
	if err != nil {
		log.Fatal(err)
	}
	rooms, err := LoadRoomStubs()
	if err != nil {
		log.Fatal(err)
	}
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	cache = Cache{
		Users:    users,
		Rooms:    rooms,
		Config:   config,
		Channels: make([]SSEStream, 0),
	}
}

func LoadConfig() (AppConfig, error) {
	adminAcc := getAdminAccount()
	var config = AppConfig{
		AdminAccount: adminAcc,
	}
	return config, nil
}
