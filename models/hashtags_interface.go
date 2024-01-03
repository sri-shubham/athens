package models

import "time"

type Hashtag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Hashtags interface {
	CRUD[*Hashtag]
}
