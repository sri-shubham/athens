package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

type PgProject struct {
	tableName   struct{}  `pg:"projects"`
	ID          int64     `pg:"id,pk"`
	Name        string    `pg:"name"`
	Slug        string    `pg:"slug"`
	Description string    `pg:"description"`
	CreatedAt   time.Time `pg:"created_at"`
}

// BeforeInsert hook is called before inserting a new record.
func (u *PgProject) BeforeInsert(ctx context.Context) (context.Context, error) {
	// Perform operations before insert
	u.CreatedAt = time.Now()

	return ctx, nil
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

// GetNewEmptyStruct implements Projects.
func (*PgProjectHelper) GetNewEmptyStruct() *Project {
	return &Project{}
}

func NewPgProjectHelper(db *pg.DB) Projects {
	return &PgProjectHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PgProject, *Project]{
			db:             db,
			MapModelToDB:   mapPgProject,
			MapModelFromDB: mapProject,
			GetEmptyStruct: func() *PgProject { return &PgProject{} },
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
