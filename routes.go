// routes.go

package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/* HTTP Request Handlers */

// GET /login
func loginPage(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", gin.H{})
}

// POST /login
func login(context *gin.Context) {
	// Read incoming email and password
	email := context.PostForm("email")
	password := context.PostForm("password")
	if email == "" || password == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	// Get a session for the user
	session := sessions.Default(context)

	// Find the user in the users list
	for _, user := range cache.Users.Users {
		// If we match an email
		if user.Email == email {
			// Check the password hash
			if CompareToHash(password, user.Hash) {
				// If the password matches, set the user's ID in the session data
				session.Set("user", user.ID.String())
				session.Save()

				// For HTMX, when using AJAX, a 200 OK status with a script to redirect is more appropriate than a 302 Found
				context.Header("Content-Type", "text/html")
				context.String(http.StatusOK, "<script>window.location.href = '/';</script>")
				return
			}
			// If password doesn't match, break out of the loop
			break
		}
	}

	// If the user was not found or the password did not match, return a 401 Unauthorized
	context.Status(http.StatusUnauthorized)
}

// GET /logout
func logout(context *gin.Context) {
	// Get a session for the user
	session := sessions.Default(context)
	// Clear the user's session
	session.Clear()
	session.Save()
	// Send a redirect script to the client to redirect to the login page
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, "<script>window.location.href = '/login';</script>")
}

// GET /register
func registerPage(context *gin.Context) {
	context.HTML(http.StatusOK, "register.html", gin.H{})
}

// POST /register
func register(context *gin.Context) {
	// Read incoming email, password, and optional display name
	email := context.PostForm("email")
	password := context.PostForm("password")
	displayName := context.DefaultPostForm("displayname", "anon") // Use "anon" as default if not provided

	if email == "" || password == "" {
		log.Println("No email or password provided")
		context.Status(http.StatusBadRequest)
		return
	}

	// Check if the email is already taken
	for _, user := range cache.Users.Users {
		if user.Email == email {
			log.Println("Email already taken")
			context.Status(http.StatusConflict)
			return
		}
	}

	// Hash the password
	hash, err := HashString(password)
	if err != nil {
		log.Println("Error hashing password")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Create a new user
	NewUser(displayName, email, hash)

	// Redirect the client to the login page
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, "<script>window.location.href = '/login';</script>")
}

// GET /
func indexPage(context *gin.Context) {
	// Send the home page
	context.HTML(http.StatusOK, "index.html", gin.H{})
}

// Get /user
func userInfo(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	user := GetUserByID(session.Get("user").(string))
	if user == nil {
		context.Status(http.StatusNotFound)
		return
	}

	userTemplate, err := UserTemplate(user)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, userTemplate)
}

// GET /rooms
// func rooms(context *gin.Context) {
// 	// Get the user's session
// 	session := sessions.Default(context)

// 	// Get the user
// 	user := GetUserByID(session.Get("user").(string))
// 	if user == nil {
// 		context.Status(http.StatusNotFound)
// 		return
// 	}

// 	// Iterate over the user's rooms ids and load the rooms, generate the rooms templates, and send the response
// 	rooms := []*Room{}
// 	for _, roomstub := range user.Rooms {
// 		room, err := LoadRoom(roomstub.ID.String())
// 		if err != nil {
// 			// Log the error
// 			log.Println(err)
// 			context.Status(http.StatusInternalServerError)
// 			return
// 		}
// 		rooms = append(rooms, &room)
// 	}

// 	roomsTemplate, err := RoomsTemplate(rooms)
// 	if err != nil {
// 		// Log the error
// 		log.Println(err)
// 		context.Status(http.StatusInternalServerError)
// 		return
// 	}

// 	context.Header("Content-Type", "text/html")
// 	context.String(http.StatusOK, roomsTemplate)
// }

// GET /room/:id
// GET /room
func room(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Get the user
	user := GetUserByID(session.Get("user").(string))
	if user == nil {
		context.Status(http.StatusNotFound)
		return
	}

	// If the user has no rooms, 404
	if len(user.Rooms) == 0 {
		context.Status(http.StatusNotFound)
		return
	}

	// Get the room ID from the URL
	roomID := context.Param("id")
	if roomID == "" {
		// Choose the users first room if no room ID is provided
		roomID = user.Rooms[0].ID.String()
	}

	// Load the room
	room, err := LoadRoom(roomID)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Generate the room template
	roomTemplate, err := RoomTemplate(&room)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, roomTemplate)
}

// POST /message/send
func sendMessage(context *gin.Context) {

	// Get the user's session
	session := sessions.Default(context)

	// Get the user
	user := GetUserByID(session.Get("user").(string))
	if user == nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Get the room ID from the form
	roomID := context.PostForm("room")
	if roomID == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	// Get the message content from the form
	content := context.PostForm("content")
	if content == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	// Create a new message
	message := NewMessage(user, content)

	// Generate the message template
	messageTemplate, err := MessageTemplate(message)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	// Get the channels for the room
	channels := GetChannels(roomID)

	// For each channel, send the message
	for _, channel := range channels {
		channel.Stream <- messageTemplate
	}

	// Load the room
	room, err := LoadRoom(roomID)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Add the message to the room
	room.Messages = append(room.Messages, *message)

	// Save the room back to the data store
	if err := SaveRoom(room); err != nil {
		// Log the error
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	// Send a 200 OK status
	context.Status(http.StatusOK)
}

// GET /tabs
func tabs(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Get the user
	user := GetUserByID(session.Get("user").(string))
	if user == nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Generate the room tabs template
	tabsTemplate, err := TabsTemplate(user.Rooms)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, tabsTemplate)
}

// GET /ping
func ping(context *gin.Context) {
	context.String(http.StatusOK, "pong")
}

// SSE GET /room/:id/stream
func streamRoom(context *gin.Context) {

	roomID := context.Param("id")
	session := sessions.Default(context)
	user := GetUserByID(session.Get("user").(string))

	context.Header("Content-Type", "text/event-stream")
	context.Header("Cache-Control", "no-cache")
	context.Header("Connection", "keep-alive")

	messageChan := GetOrCreateChannel(roomID, user.ID.String())

	clientGone := context.Request.Context().Done()

	for {
		select {
		case message := <-messageChan.Stream:
			context.SSEvent("message", message)
			if flusher, ok := context.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-clientGone:
			log.Println("Client disconnected")
			RemoveChannel(roomID, user.ID.String())
			return
		}
	}
}
