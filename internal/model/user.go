package model

import "time"

type User struct {
	ID        int       `json:"ID,omitempty"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
