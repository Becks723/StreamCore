package model

import "time"

// VideoEs - feed es database
type VideoEs struct {
	Id          uint      `json:"id"`
	AuthorId    uint      `json:"author_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Username    string    `json:"username"`
}
