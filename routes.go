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
			if compareToHash(password, user.Hash) {
				// If the password matches, set the user's ID in the session data
				session.Set("user", user.ID.String())
				session.Save()

				// For HTMX, a 200 OK status with a script to redirect is more appropriate than a 302 Found
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
	hash, err := hashString(password)
	if err != nil {
		log.Println("Error hashing password")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Add the new user to the users list
	newUser := NewUser(displayName, email, hash)
	cache.Users.Users = append(cache.Users.Users, *newUser)

	// Save the updated users back to the data store
	if err := SaveUsers(cache.Users); err != nil {
		log.Println("Error saving users")
		context.Status(http.StatusInternalServerError)
		return
	}

	// Redirect the client to the login page
	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, "<script>window.location.href = '/login';</script>")
}

// GET /
func indexPage(context *gin.Context) {
	// Check if we have a session and if not redirect to the login page
	session := sessions.Default(context)
	if session.Get("user") == nil {
		context.Redirect(http.StatusFound, "/login")
		return
	}

	// Send the home page
	context.HTML(http.StatusOK, "index.html", gin.H{})
}

// Get /user
func userInfo(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if we have a session and if not redirect to the login page
	if session.Get("user") == nil {
		context.Redirect(http.StatusFound, "/login")
		return
	}

	user := getUserByID(session.Get("user").(string))
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
func rooms(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if we have a session and if not return 404
	if session.Get("user") == nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Get the user
	user := getUserByID(session.Get("user").(string))
	if user == nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Iterate over the user's rooms ids and load the rooms, generate the rooms templates, and send the response
	rooms := []*Room{}
	for _, roomID := range user.Rooms {
		room, err := LoadRoom(roomID.String())
		if err != nil {
			// Log the error
			log.Println(err)
			context.Status(http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, &room)
	}

	roomsTemplate, err := RoomsTemplate(rooms)
	if err != nil {
		// Log the error
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, roomsTemplate)
}

// POST /message/send
func sendMessage(context *gin.Context) {
	// Get the user's session
	session := sessions.Default(context)

	// Check if we have a session and if not return 404
	if session.Get("user") == nil {
		context.Status(http.StatusNotFound)
		return
	}

	// Get the user
	user := getUserByID(session.Get("user").(string))
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

	// Load the room
	room, err := LoadRoom(roomID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
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

	// Add the message to the room
	room.Messages = append(room.Messages, *message)

	// Save the room back to the data store
	if err := SaveRoom(room); err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	// Generate the message template
	messageTemplate, err := MessageTemplate(message)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Header("Content-Type", "text/html")
	context.String(http.StatusOK, messageTemplate)
}

// GET /ping
func ping(context *gin.Context) {
	context.String(http.StatusOK, "pong")
}
