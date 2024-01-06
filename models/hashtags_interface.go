package models

import "time"

type Hashtag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Hashtag) GetID() int64 {
	if h == nil {
		return 0
	}
	return h.ID
}

type Hashtags interface {
	CRUD[*Hashtag]
}
