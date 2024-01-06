package api

import "github.com/sri-shubham/athens/models"

var _ = (GenericCrud)(&UsersCrud{})

type ProjectHashtagsCrud struct {
	GenericCrudHelper[*models.ProjectHashtag]
}

func NewProjectHashtagsCrud(db models.ProjectHashtags) GenericCrud {
	return &ProjectHashtagsCrud{
		GenericCrudHelper: GenericCrudHelper[*models.ProjectHashtag]{
			db: db,
		},
	}
}
