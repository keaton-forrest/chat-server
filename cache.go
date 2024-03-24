// cache.go

package main

import "log"

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
	cache = Cache{
		Users: users,
		Rooms: rooms,
	}
}
