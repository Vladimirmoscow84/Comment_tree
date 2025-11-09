package model

import "time"

type Comment struct {
	ID        int       `json:"id"`
	ParentID  int       `json:"parent_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
