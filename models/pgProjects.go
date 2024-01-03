package models

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type PgProject struct {
	ID          int64     `pg:"id,pk"`
	Name        string    `pg:"name"`
	Slug        string    `pg:"slug"`
	Description string    `pg:"description"`
	CreatedAt   time.Time `pg:"created_at"`
}

func (p *PgProject) GetID() int64 {
	if p == nil {
		return 0
	}
	return p.ID
}

// Checks interface is implemented
var _ = Projects(&PgProjectHelper{})

// PgUser: Postgres implementation of users_interface
type PgProjectHelper struct {
	db *pg.DB
	*CRUDHelper[*PgProject, *Project]
}

func NewPgProjectHelper(db *pg.DB) Projects {
	return &PgProjectHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgProject, *Project]{
			db:             db,
			MapModelToDB:   mapPgProject,
			MapModelFromDB: mapProject,
		},
	}
}

func mapProject(in *PgProject) *Project {
	return &Project{
		ID:          in.ID,
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		CreatedAt:   in.CreatedAt,
	}
}

func mapPgProject(in *Project) *PgProject {
	return &PgProject{
		ID:          in.ID,
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		CreatedAt:   in.CreatedAt,
	}
}
