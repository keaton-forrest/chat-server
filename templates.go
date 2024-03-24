// templates.go

package main

import (
	"bytes"
	"html/template"
)

/* Templates */

// HTML templates for the chat application
//  - Rooms
// 		- Room
// 			- Message history
// 				- Message
// 			- Users
// 				- User
//  		- Input

// Define a global variable for templates. This is compiled once at startup.
var templates *template.Template

func init() {
	// Initialize and parse all templates.
	templates = template.Must(template.New("master").Parse(`
{{define "rooms"}}
	{{range .}}
		{{template "room" .}}
	{{end}}
{{end}}

{{define "room"}}
<div class='room' id='id-{{.ID}}'>
	<h2>{{.Name}}</h2>
	<div>
		<div class='messages' hx-ext="sse" sse-connect='/room/{{.ID}}' sse-swap='message' hx-swap='beforeend'>
			{{range .Messages}}
				{{template "message" .}}
			{{end}}
		</div>
		<div class='users'>
			{{template "users" .Mods}}
			{{template "users" .Users}}
		</div>
	</div>
	{{template "input" .}}
</div>
{{end}}

{{define "message"}}
<div class='message {{.Status}}'>
	<div class='author'>{{.Author.DisplayName}}</div>
	<div class='content'>{{.Content}}</div>
	<div class='meta'>
		<span class='created'>{{.CreatedAt}}</span>
		{{if ne .ModifiedAt ""}}
			<span class='modified'>Modified: {{.ModifiedAt}}</span>
		{{end}}
	</div>
</div>
{{end}}

{{define "users"}}
	{{range .}}
		{{template "user" .}}
	{{end}}
{{end}}

{{define "user"}}
<div class='user'>
	<div class='displayname'>{{.DisplayName}}</div>
</div>
{{end}}

{{define "input"}}
<div class='input'>
	<form hx-post='/message/send' hx-indicator="#loading" hx-swap='none'>
		<input type='text' name='content' required>
		<input type='hidden' name='room' value='{{.ID}}'>
		<button type='submit'>Send</button>
	</form>
</div>
{{end}}
`))
}

// RoomsTemplate generates HTML for rooms.
func RoomsTemplate(rooms []*Room) (string, error) {
	var htmlResponse bytes.Buffer
	if err := templates.ExecuteTemplate(&htmlResponse, "rooms", rooms); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

func UserTemplate(user *User) (string, error) {
	var htmlResponse bytes.Buffer
	if err := templates.ExecuteTemplate(&htmlResponse, "user", user); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

func UsersTemplate(users []*User) (string, error) {
	var htmlResponse bytes.Buffer
	if err := templates.ExecuteTemplate(&htmlResponse, "users", users); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

func MessageTemplate(message *Message) (string, error) {
	var htmlResponse bytes.Buffer
	if err := templates.ExecuteTemplate(&htmlResponse, "message", message); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}
