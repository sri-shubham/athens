package models

import "time"

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}

type Users interface {
	CRUD[*User]
}