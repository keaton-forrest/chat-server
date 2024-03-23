// templates.go

package main

import (
	"bytes"
	"html/template"
)

/* Templates */

// HTML templates for the chat application
// 	- Room
// 		- Message history
// 			- Message
// 		- Users
// 			- User
//  	- Input

// Room template
func RoomTemplate(room Room) (string, error) {
	tmpl := template.Must(template.New("room").Parse(`
		<div class='room'>
			<div class='messages'>
				{{range .Messages}}
					{{MessageTemplate .}}
				{{end}}
			</div>
			<div class='users'>
				{{range .Users}}
					{{UserTemplate .}}
				{{end}}
			</div>
		</div>
	`))

	var htmlResponse bytes.Buffer
	if err := tmpl.Execute(&htmlResponse, room); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

// Message template
func MessageTemplate(message Message) (string, error) {
	tmpl := template.Must(template.New("message").Parse(`
		<div class='message'>
			<div class='author'>{{.Author.DisplayName}}</div>
			<div class='content'>{{.Content}}</div>
		</div>
	`))

	var htmlResponse bytes.Buffer
	if err := tmpl.Execute(&htmlResponse, message); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

// User template
func UserTemplate(user User) (string, error) {
	tmpl := template.Must(template.New("user").Parse(`
		<div class='user'>
			<div class='displayname'>{{.DisplayName}}</div>
		</div>
	`))

	var htmlResponse bytes.Buffer
	if err := tmpl.Execute(&htmlResponse, user); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}

// Input template
func InputTemplate() (string, error) {
	tmpl := template.Must(template.New("input").Parse(`
		<div class='input'>
			<form hx-post='/message/send' hx-swap='outerHTML' hx-indicator="#loading">
				<input type='text' name='content' required>
				<button type='submit'>Send</button>
			</form>
		</div>
	`))

	var htmlResponse bytes.Buffer
	if err := tmpl.Execute(&htmlResponse, nil); err != nil {
		return "", err
	}
	return htmlResponse.String(), nil
}
