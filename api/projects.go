package api

import "github.com/sri-shubham/athens/models"

var _ = (GenericCrud)(&UsersCrud{})

type ProjectsCrud struct {
	GenericCrudHelper[*models.Project]
}

func NewProjectsCrud(db models.Projects) GenericCrud {
	return &ProjectsCrud{
		GenericCrudHelper: GenericCrudHelper[*models.Project]{
			db: db,
		},
	}
}
