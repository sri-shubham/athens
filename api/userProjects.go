package api

import "github.com/sri-shubham/athens/models"

var _ = (GenericCrud)(&UsersCrud{})

type UserProjectsCrud struct {
	GenericCrudHelper[*models.UserProject]
}

func NewUserProjectsCrud(db models.UserProjects) GenericCrud {
	return &UserProjectsCrud{
		GenericCrudHelper: GenericCrudHelper[*models.UserProject]{
			db: db,
		},
	}
}
