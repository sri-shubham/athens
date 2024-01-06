package models

import "context"

type UserProject struct {
	ID        int64 `json:"id"`
	ProjectID int64 `json:"project_id"`
	UserID    int64 `json:"user_id"`
}

func (u *UserProject) GetID() int64 {
	if u == nil {
		return 0
	}
	return u.ID
}

type UserProjects interface {
	CRUD[*UserProject]
	GetByProjectID(ctx context.Context, id int64) ([]*UserProject, error)
}
