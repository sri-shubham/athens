package models

import (
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

func NewPgUserHelper(db *pg.DB) Users {
	return &PgUserHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgUser, *User]{
			db:             db,
			MapModelToDB:   mapPgUser,
			MapModelFromDB: mapUser,
		},
	}
}

type PgUser struct {
	ID        int64  `pg:",pk"`
	Name      string `pg:",notnull"`
	CreatedAt *time.Time
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
