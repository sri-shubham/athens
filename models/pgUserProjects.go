package models

import (
	"github.com/go-pg/pg/v10"
)

type PgUserProjectHelper struct {
	*CRUDHelper[*PGUserProject, *UserProject]
	db *pg.DB
}

// GetNewEmptyStruct implements UserProjects.
func (*PgUserProjectHelper) GetNewEmptyStruct() *UserProject {
	return &UserProject{}
}

func NewPgUserProjectHelper(db *pg.DB) UserProjects {
	return &PgUserProjectHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PGUserProject, *UserProject]{
			db:             db,
			MapModelToDB:   mapPgUserProject,
			MapModelFromDB: mapUserProject,
			GetEmptyStruct: func() *PGUserProject { return &PGUserProject{} },
		},
	}
}

var _ = UserProjects(&PgUserProjectHelper{})

type PGUserProject struct {
	tableName struct{} `pg:"user_projects"`
	ID        int64    `pg:"id,pk"`
	ProjectID int64    `pg:"project_id"`
	UserID    int64    `pg:"user_id"`
}

func (p *PGUserProject) GetID() int64 {
	if p == nil {
		return 0
	}
	return p.ID
}

func mapUserProject(in *PGUserProject) *UserProject {
	return &UserProject{
		ID:        in.ID,
		ProjectID: in.ProjectID,
		UserID:    in.UserID,
	}
}

func mapPgUserProject(in *UserProject) *PGUserProject {
	return &PGUserProject{
		ID:        in.ID,
		ProjectID: in.ProjectID,
		UserID:    in.UserID,
	}
}
