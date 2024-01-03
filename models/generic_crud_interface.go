package models

import "github.com/go-pg/pg/v10"

// Generic CRUD operations interface
type CRUD[T any] interface {
	Get(id int64) (T, error)
	Create(T) (int64, error)
	Update(T) (int64, error)
	Delete(id int64) error
}

type IdGetter interface {
	GetID() int64
}

// We are adding generic Implementation so we can
// to reduce code duplication
type CRUDHelper[DbType IdGetter, ModelType any] struct {
	db             *pg.DB
	MapModelToDB   func(ModelType) DbType
	MapModelFromDB func(DbType) ModelType
}

func (uh *CRUDHelper[DbType, ModelType]) Get(id int64) (ModelType, error) {
	var dbItem DbType
	err := uh.db.Model(dbItem).WherePK().Select()
	if err != nil {
		var v ModelType
		return v, err
	}
	return uh.MapModelFromDB(dbItem), nil
}

func (uh *CRUDHelper[DbType, ModelType]) Create(u ModelType) (int64, error) {
	dbItem := uh.MapModelToDB(u)
	_, err := uh.db.Model(dbItem).Insert()
	if err != nil {
		return 0, err
	}
	return dbItem.GetID(), nil
}

func (uh *CRUDHelper[DbType, ModelType]) Delete(id int64) error {
	var pgModel DbType
	_, err := uh.db.Model(pgModel).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}

func (uh *CRUDHelper[DbType, ModelType]) Update(in ModelType) (int64, error) {
	dbItem := uh.MapModelToDB(in)
	_, err := uh.db.Model(dbItem).WherePK().Update()
	if err != nil {
		return 0, err
	}
	return dbItem.GetID(), nil
}