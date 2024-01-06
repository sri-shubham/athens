package models

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sri-shubham/athens/util"
)

type PgUserProjectHelper struct {
	*CRUDHelper[*PGUserProject, *UserProject]
	db *pg.DB
}

// GetNewEmptyStruct implements UserProjects.
func (*PgUserProjectHelper) GetNewEmptyStruct() *UserProject {
	return &UserProject{}
}

func NewPgUserProjectHelper(db *pg.DB, updateQueue util.UpdateQueue) UserProjects {
	return &PgUserProjectHelper{
		db: db,
		CRUDHelper: &CRUDHelper[*PGUserProject, *UserProject]{
			db:             db,
			updateQueue:    updateQueue,
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

// GetByProjectID implements ProjectHashtags.
func (h *PgUserProjectHelper) GetByProjectID(ctx context.Context, id int64) ([]*UserProject, error) {
	var dbItem []*PGUserProject
	var out []*UserProject
	err := h.db.Model(&dbItem).Context(ctx).Where("project_id=?", id).Select()
	if err != nil {
		return out, err
	}

	out = make([]*UserProject, 0, len(dbItem))
	for _, item := range dbItem {
		out = append(out, h.MapModelFromDB(item))
	}

	return out, nil
}
