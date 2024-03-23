// utils.go

package main

import "time"

/* Utility Functions */

// currentTime returns the current time as a string in RFC 3339 format.
func currentTime() string {
	return time.Now().Format(time.RFC3339)
}
