package models

import "context"

type ProjectHashtag struct {
	ID        int64 `json:"id"`
	HashtagID int64 `json:"hashtag_id"`
	ProjectID int64 `json:"project_id"`
}

func (u *ProjectHashtag) GetID() int64 {
	if u == nil {
		return 0
	}
	return u.ID
}

type ProjectHashtags interface {
	CRUD[*ProjectHashtag]
	GetByHashTag(ctx context.Context, id int64) ([]*ProjectHashtag, error)
	GetByProjectID(ctx context.Context, id int64) ([]*ProjectHashtag, error)
}
