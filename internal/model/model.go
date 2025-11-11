package model

import "time"

type Comment struct {
	ID        int       `json:"id" db:"id"`
	ParentID  *int      `json:"parent_id" db:"parent_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_ar"`
}
