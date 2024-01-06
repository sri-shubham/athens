package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

// Checks interface is implemented
var _ = Users(&PgUserHelper{})

// PgUser: Postgres implementation of users_interface
type PgUserHelper struct {
	db *pg.DB
	*CRUDHelper[*PgUser, *User]
}

// GetNewEmptyStruct implements Users.
func (*PgUserHelper) GetNewEmptyStruct() *User {
	panic("unimplemented")
}

func NewPgUserHelper(db *pg.DB) Users {
	return &PgUserHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgUser, *User]{
			db:             db,
			MapModelToDB:   mapPgUser,
			MapModelFromDB: mapUser,
			GetEmptyStruct: func() *PgUser { return &PgUser{} },
		},
	}
}

type PgUser struct {
	tableName struct{} `pg:"users"`
	ID        int64    `pg:",pk"`
	Name      string   `pg:",notnull"`
	CreatedAt time.Time
}

// BeforeInsert hook is called before inserting a new record.
func (u *PgUser) BeforeInsert(ctx context.Context) (context.Context, error) {
	// Perform operations before insert
	u.CreatedAt = time.Now()

	return ctx, nil
}

func (p *PgUser) GetID() int64 {
	if p == nil {
		return 0
	}
	return p.ID
}

func mapUser(in *PgUser) *User {
	return &User{
		ID:        in.ID,
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
	}
}

func mapPgUser(in *User) *PgUser {
	return &PgUser{
		ID:        in.ID,
		Name:      in.Name,
		CreatedAt: in.CreatedAt,
	}
}
