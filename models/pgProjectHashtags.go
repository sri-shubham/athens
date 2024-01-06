package models

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sri-shubham/athens/util"
)

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

// GetByHashTag implements ProjectHashtags.
func (h *PgProjectHashtagHelper) GetByHashTag(ctx context.Context, id int64) ([]*ProjectHashtag, error) {
	var dbItem []*PgProjectHashtag
	var out []*ProjectHashtag
	err := h.db.Model(&dbItem).Context(ctx).Where("hashtag_id=?", id).Select()
	if err != nil {
		return out, err
	}

	out = make([]*ProjectHashtag, 0, len(dbItem))
	for _, item := range dbItem {
		out = append(out, h.MapModelFromDB(item))
	}

	return out, nil
}

// GetByProjectID implements ProjectHashtags.
func (h *PgProjectHashtagHelper) GetByProjectID(ctx context.Context, id int64) ([]*ProjectHashtag, error) {
	var dbItem []*PgProjectHashtag
	var out []*ProjectHashtag
	err := h.db.Model(&dbItem).Context(ctx).Where("project_id=?", id).Select()
	if err != nil {
		return out, err
	}

	out = make([]*ProjectHashtag, 0, len(dbItem))
	for _, item := range dbItem {
		out = append(out, h.MapModelFromDB(item))
	}

	return out, nil
}

// GetNewEmptyStruct implements ProjectHashtags.
func (*PgProjectHashtagHelper) GetNewEmptyStruct() *ProjectHashtag {
	return &ProjectHashtag{}
}

func NewPgProjectHashtagHelper(db *pg.DB, updateQueue util.UpdateQueue) ProjectHashtags {
	return &PgProjectHashtagHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgProjectHashtag, *ProjectHashtag]{
			db:             db,
			updateQueue:    updateQueue,
			MapModelToDB:   mapPgProjectHashtag,
			MapModelFromDB: mapProjectHashtag,
			GetEmptyStruct: func() *PgProjectHashtag { return &PgProjectHashtag{} },
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
