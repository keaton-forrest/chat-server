// types.go

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/* Types */

type Cache struct {
	Users    Users
	Rooms    Roomstubs
	Config   AppConfig
	Channels []SSEStream
}

type SSEStream struct {
	Stream chan string
	RoomID uuid.UUID
	UserID uuid.UUID
}

type AppConfig struct {
	AdminAccount gin.Accounts
}

// Message Status types
const (
	StatusNew     = "new"
	StatusEdited  = "edited"
	StatusDeleted = "deleted"
)

type Message struct {
	ID         uuid.UUID `json:"id"`
	Author     *UserStub `json:"author"`
	Content    string    `json:"content"`
	CreatedAt  string    `json:"created_at"`
	ModifiedAt string    `json:"modified_at"`
	Status     string    `json:"status"`
}

type Room struct {
	Name       string     `json:"name"`
	ID         uuid.UUID  `json:"id"`
	Messages   []Message  `json:"messages"`
	Mods       []UserStub `json:"mods"`
	Users      []UserStub `json:"users"`
	Banned     []UserStub `json:"banned"`
	Muted      []UserStub `json:"muted"`
	CreatedAt  string     `json:"created_at"`
	ModifiedAt string     `json:"modified_at"`
}

type User struct {
	ID          uuid.UUID   `json:"id"`
	DisplayName string      `json:"displayname"`
	Email       string      `json:"email"`
	Hash        string      `json:"hash"`
	Rooms       []*RoomStub `json:"rooms"`
	CreatedAt   string      `json:"created_at"`
	ModifiedAt  string      `json:"modified_at"`
}

type Users struct {
	Users []User `json:"users"`
}

type UserStub struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"displayname"`
}

type RoomStub struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Roomstubs struct {
	Rooms []RoomStub `json:"rooms"`
}

/* Methods for Message */

// NewMessage creates a new Message instance with initial values.
func NewMessage(author *User, content string) *Message {
	return &Message{
		ID:         uuid.New(),
		Author:     &UserStub{ID: author.ID, DisplayName: author.DisplayName},
		Content:    content,
		CreatedAt:  CurrentTime(),
		ModifiedAt: CurrentTime(),
		Status:     StatusNew,
	}
}

// UpdateModifiedAt updates the ModifiedAt timestamp of the Message.
func (m *Message) UpdateModifiedAt() {
	m.ModifiedAt = CurrentTime()
	m.Status = StatusEdited
}

// Delete marks the Message as deleted.
func (m *Message) Delete() {
	m.ModifiedAt = CurrentTime()
	m.Status = StatusDeleted
}

/* Methods for Room */

// NewRoom creates a new Room instance with initial values.
func NewRoom(name string) *Room {
	room := &Room{
		Name:       name,
		ID:         uuid.New(),
		Messages:   []Message{},
		Mods:       []UserStub{},
		Users:      []UserStub{},
		Banned:     []UserStub{},
		Muted:      []UserStub{},
		CreatedAt:  CurrentTime(),
		ModifiedAt: CurrentTime(),
	}
	roomStub := RoomStub{
		ID:   room.ID,
		Name: room.Name,
	}
	cache.Rooms.Rooms = append(cache.Rooms.Rooms, roomStub)
	// Save the rooms to file
	SaveRoomStubs(cache.Rooms)
	return room
}

// UpdateModifiedAt updates the ModifiedAt timestamp of the Room.
func (r *Room) UpdateModifiedAt() {
	r.ModifiedAt = CurrentTime()
}

/* Methods for User */

// NewUser creates a new User instance with initial values.
func NewUser(displayName, email, hash string) *User {
	user := User{
		ID:          uuid.New(),
		DisplayName: displayName,
		Email:       email,
		Hash:        hash,
		// Add the default room ID for General Chat
		// {
		//     "id": "260ca734-06ff-4baa-baaf-8e440730e960",
		//     "name": "General Chat"
		// }
		Rooms:      []*RoomStub{GetRoomStubByID("260ca734-06ff-4baa-baaf-8e440730e960")},
		CreatedAt:  CurrentTime(),
		ModifiedAt: CurrentTime(),
	}
	cache.Users.Users = append(cache.Users.Users, user)
	// Save the users to file
	SaveUsers(cache.Users)
	return &user
}

// UpdateModifiedAt updates the ModifiedAt timestamp of the User.
func (u *User) UpdateModifiedAt() {
	u.ModifiedAt = CurrentTime()
}
