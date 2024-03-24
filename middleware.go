// middleware.go

package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionAuthMiddleware checks if the user is authenticated via session.
// If not authenticated, it returns a 401 Unauthorized response.
func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("authenticated")
		if user == nil {
			// User is not logged in, return a 401 Unauthorized response
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// Proceed to the next handler if the user is authenticated
		c.Next()
	}
}
