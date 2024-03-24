// types.go

package main

import (
	"github.com/google/uuid"
)

/* Types */

// Message Status types
const (
	StatusNew     = "new"
	StatusEdited  = "edited"
	StatusDeleted = "deleted"
)

type Message struct {
	ID         uuid.UUID `json:"id"`
	Author     *User     `json:"author"`
	Content    string    `json:"content"`
	CreatedAt  string    `json:"created_at"`
	ModifiedAt string    `json:"modified_at"`
	Status     string    `json:"status"`
}

type Room struct {
	Name       string      `json:"name"`
	ID         uuid.UUID   `json:"id"`
	Messages   []Message   `json:"messages"`
	Mods       []uuid.UUID `json:"mods"`
	Users      []uuid.UUID `json:"users"`
	CreatedAt  string      `json:"created_at"`
	ModifiedAt string      `json:"modified_at"`
}

type User struct {
	ID          uuid.UUID   `json:"id"`
	DisplayName string      `json:"displayname"`
	Email       string      `json:"email"`
	Hash        string      `json:"hash"`
	Rooms       []uuid.UUID `json:"rooms"`
	CreatedAt   string      `json:"created_at"`
	ModifiedAt  string      `json:"modified_at"`
}

type Users struct {
	Users []User `json:"users"`
}

/* Methods for Message */

// NewMessage creates a new Message instance with initial values.
func NewMessage(author *User, content string) *Message {
	return &Message{
		ID:         uuid.New(),
		Author:     author,
		Content:    content,
		CreatedAt:  currentTime(),
		ModifiedAt: currentTime(),
		Status:     StatusNew,
	}
}

// UpdateModifiedAt updates the ModifiedAt timestamp of the Message.
func (m *Message) UpdateModifiedAt() {
	m.ModifiedAt = currentTime()
	m.Status = StatusEdited
}

// Delete marks the Message as deleted.
func (m *Message) Delete() {
	m.ModifiedAt = currentTime()
	m.Status = StatusDeleted
}

/* Methods for Room */

// NewRoom creates a new Room instance with initial values.
func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		ID:         uuid.New(),
		CreatedAt:  currentTime(),
		ModifiedAt: currentTime(),
	}
}

// UpdateModifiedAt updates the ModifiedAt timestamp of the Room.
func (r *Room) UpdateModifiedAt() {
	r.ModifiedAt = currentTime()
}

/* Methods for User */

// NewUser creates a new User instance with initial values.
func NewUser(displayName, email, hash string) *User {
	return &User{
		ID:          uuid.New(),
		DisplayName: displayName,
		Email:       email,
		Hash:        hash,
		// Add the default room ID for General Chat
		Rooms:      []uuid.UUID{uuid.MustParse("260ca734-06ff-4baa-baaf-8e440730e960")},
		CreatedAt:  currentTime(),
		ModifiedAt: currentTime(),
	}
}

// UpdateModifiedAt updates the ModifiedAt timestamp of the User.
func (u *User) UpdateModifiedAt() {
	u.ModifiedAt = currentTime()
}

// GetUserByID returns a pointer to the User with the given ID.
func getUserByID(users []User, id string) *User {
	for _, user := range users {
		if user.ID == uuid.MustParse(id) {
			return &user
		}
	}
	return nil
}
