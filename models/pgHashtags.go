package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

type PgHashtag struct {
	tableName struct{}  `pg:"hashtags"`
	ID        int64     `pg:"id,pk"`
	Name      string    `pg:"name"`
	CreatedAt time.Time `pg:"created_at"`
}

// BeforeInsert hook is called before inserting a new record.
func (u *PgHashtag) BeforeInsert(ctx context.Context) (context.Context, error) {
	// Perform operations before insert
	u.CreatedAt = time.Now()

	return ctx, nil
}

// Checks interface is implemented
var _ = Users(&PgUserHelper{})

// PgUser: Postgres implementation of users_interface
type PgHashtagHelper struct {
	db *pg.DB
	*CRUDHelper[*PgHashtag, *Hashtag]
}

// GetNewEmptyStruct implements Hashtags.
func (*PgHashtagHelper) GetNewEmptyStruct() *Hashtag {
	return &Hashtag{}
}

func NewPgHashtagHelper(db *pg.DB) Hashtags {
	return &PgHashtagHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgHashtag, *Hashtag]{
			db:             db,
			MapModelToDB:   mapPgHashtag,
			MapModelFromDB: mapHashtag,
			GetEmptyStruct: func() *PgHashtag { return &PgHashtag{} },
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
