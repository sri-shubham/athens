package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type PgHashtag struct {
	ID        int64     `pg:"id,pk"`
	Name      string    `pg:"name"`
	CreatedAt time.Time `pg:"created_at"`
}

// Checks interface is implemented
var _ = Users(&PgUserHelper{})

// PgUser: Postgres implementation of users_interface
type PgHashtagHelper struct {
	db *pg.DB
	*CRUDHelper[*PgHashtag, *Hashtag]
}

func NewPgHashtagHelper(db *pg.DB) Hashtags {
	return &PgHashtagHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgHashtag, *Hashtag]{
			db:             db,
			MapModelToDB:   mapPgHashtag,
			MapModelFromDB: mapHashtag,
		},
	}
}

func (p *PgHashtag) GetID() int64 {
	if p == nil {
		return 0
	}
	return p.ID
}

func mapHashtag(in *PgHashtag) *Hashtag {
	return &Hashtag{
		ID:        in.ID,
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
	}
}

func mapPgHashtag(in *Hashtag) *PgHashtag {
	return &PgHashtag{
		ID:        in.ID,
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
	}
}
