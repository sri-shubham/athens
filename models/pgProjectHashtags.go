package models

import "github.com/go-pg/pg/v10"

type PgProjectHashtag struct {
	tableName struct{} `pg:"project_hashtags"`
	ID        int64    `pg:"id,pk"`
	HashtagID int64    `pg:"hashtag_id"`
	ProjectID int64    `pg:"project_id"`
}

func (p *PgProjectHashtag) GetID() int64 {
	if p == nil {
		return 0
	}
	return p.ID
}

// Checks interface is implemented
var _ = ProjectHashtags(&PgProjectHashtagHelper{})

// PgUser: Postgres implementation of users_interface
type PgProjectHashtagHelper struct {
	db *pg.DB
	*CRUDHelper[*PgProjectHashtag, *ProjectHashtag]
}

func NewPgProjectHashtagHelper(db *pg.DB) ProjectHashtags {
	return &PgProjectHashtagHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgProjectHashtag, *ProjectHashtag]{
			db:             db,
			MapModelToDB:   mapPgProjectHashtag,
			MapModelFromDB: mapProjectHashtag,
		},
	}
}

func mapProjectHashtag(in *PgProjectHashtag) *ProjectHashtag {
	return &ProjectHashtag{
		ID:        in.ID,
		HashtagID: in.HashtagID,
		ProjectID: in.ProjectID,
	}
}

func mapPgProjectHashtag(in *ProjectHashtag) *PgProjectHashtag {
	return &PgProjectHashtag{
		ID:        in.ID,
		HashtagID: in.HashtagID,
		ProjectID: in.ProjectID,
	}
}
