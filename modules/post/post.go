package post

import (
	"time"

	"github.com/kielkow/Post-Service/modules/author"
)

// Post struct
type Post struct {
	ID          int           `json:"id"`
	Author      author.Author `json:"author"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// CreatePost struct
type CreatePost struct {
	AuthorID    int    `json:"authorId"`
	Description string `json:"description"`
}

// UpdatePost struct
type UpdatePost struct {
	AuthorID    int    `json:"authorId"`
	Description string `json:"description"`
}

