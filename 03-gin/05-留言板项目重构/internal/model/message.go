package model

import "time"

type Message struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname" form:"nickname" binding:"required"`
	Content   string    `json:"content" form:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
