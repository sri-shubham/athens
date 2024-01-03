package models

type UserProject struct {
	ID        int64 `json:"id"`
	ProjectID int64 `json:"project_id"`
	UserID    int64 `json:"user_id"`
}

type UserProjects interface {
	CRUD[*UserProject]
}
