// router.go

package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes the HTTP request handlers and starts the HTTP server
func InitRouter() {

	// Create a new Gin router
	router := gin.Default()

	// Create a new cookie store
	store := cookie.NewStore([]byte("93wt4nvp2y34b223r2wet43"))

	// Use the cookie store for the session middleware
	router.Use(sessions.Sessions("session", store))

	// Load HTML files
	router.LoadHTMLGlob("templates/*")

	// Serve static files
	router.Static("/static", "./static")

	// Register the HTTP request handlers
	RegisterHandlers(router)

	// Start the HTTP server on port 8080
	router.Run(":8080")

}

// HTTP Request Handlers
func RegisterHandlers(router *gin.Engine) {

	// Basic auth middleware
	isAdmin := router.Group("/", gin.BasicAuth(cache.Config.AdminAccount))

	// Session auth middleware
	isAuthenticated := router.Group("/", SessionAuthMiddleware())

	// Register
	isAdmin.GET("/register", registerPage)
	isAdmin.POST("/register", register)
	// Login
	router.GET("/login", loginPage)
	router.POST("/login", login)
	// Logout
	router.GET("/logout", logout)
	// Index
	isAuthenticated.GET("/", indexPage)
	isAuthenticated.GET("/user", userInfo)
	isAuthenticated.GET("/room/:id", room)
	isAuthenticated.GET("/room", room)
	isAuthenticated.GET("/tabs", tabs)
	isAuthenticated.POST("/message/send", sendMessage)
	isAuthenticated.GET("/room/:id/stream", streamRoom)
	// Health check
	router.GET("/ping", ping)
}
