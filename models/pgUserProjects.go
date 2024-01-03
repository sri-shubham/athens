package models

import "github.com/go-pg/pg/v10"

type PgUserProjectHelper struct {
	*CRUDHelper[*PGUserProject, *UserProject]
	db *pg.DB
}

func NewPgUserProjectHelper(db *pg.DB) UserProjects {
	return &PgUserProjectHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PGUserProject, *UserProject]{
			db:             db,
			MapModelToDB:   mapPgUserProject,
			MapModelFromDB: mapUserProject,
		},
	}
}

var _ = UserProjects(&PgUserProjectHelper{})

type PGUserProject struct {
	ID        int64 `pg:"id,pk"`
	ProjectID int64 `pg:"project_id"`
	UserID    int64 `pg:"user_id"`
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