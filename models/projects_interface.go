package models

import "time"

type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *Project) GetID() int64 {
	if p == nil {
		return 0
	}
	return p.ID
}

type Projects interface {
	CRUD[*Project]
}
