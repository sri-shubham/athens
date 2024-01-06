package api

import "github.com/sri-shubham/athens/models"

var _ = (GenericCrud)(&UsersCrud{})

type HashtagsCrud struct {
	GenericCrudHelper[*models.Hashtag]
}

func NewHashtagsCrud(db models.Hashtags) GenericCrud {
	return &HashtagsCrud{
		GenericCrudHelper: GenericCrudHelper[*models.Hashtag]{
			db: db,
		},
	}
}
