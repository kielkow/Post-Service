package post

import "time"

// Post struct
type Post struct {
	id          int       `json:"id"`
	author      string    `json:"author"`
	description string    `json:"description"`
	createdAt   time.Time `json:"createdAt"`
	updatedAt   time.Time `json:"updatedAt"`
}
