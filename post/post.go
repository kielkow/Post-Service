package post

import "time"

// Post struct
type Post struct {
	ID          int    `json:"id"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreatePost struct
type CreatePost struct {
	Author      string `json:"author"`
	Description string `json:"description"`
}
