package api

import "github.com/sri-shubham/athens/models"

var _ = (GenericCrud)(&UsersCrud{})

type UsersCrud struct {
	GenericCrudHelper[*models.User]
}

func NewUsersCrud(db models.Users) GenericCrud {
	return &UsersCrud{
		GenericCrudHelper: GenericCrudHelper[*models.User]{
			db: db,
		},
	}
}
