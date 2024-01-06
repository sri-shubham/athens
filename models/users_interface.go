package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) GetID() int64 {
	if u == nil {
		return 0
	}
	return u.ID
}

type Users interface {
	CRUD[*User]
}
