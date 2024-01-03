package models

type ProjectHashtag struct {
	ID        int64 `json:"id"`
	HashtagID int64 `json:"hashtag_id"`
	ProjectID int64 `json:"project_id"`
}

type ProjectHashtags interface {
	CRUD[*ProjectHashtag]
}
