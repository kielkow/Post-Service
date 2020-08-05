package author

import (
	"time"
)

// Author struct
type Author struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Avatar    *string   `json:"avatar"`
	AvatarURL *string   `json:"avatarUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateAuthor struct
type CreateAuthor struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// CreateAuthorAvatar struct
type CreateAuthorAvatar struct {
	AuthorID int    `json:"authorId"`
	Avatar   string `json:"avatar"`
}
