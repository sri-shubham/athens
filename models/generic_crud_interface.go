package models

import "github.com/go-pg/pg/v10"

// Generic CRUD operations interface
type CRUD[T any] interface {
	GetNewEmptyStruct() T
	Get(id int64) (T, error)
	Create(T) (int64, error)
	Update(T) (int64, error)
	Delete(id int64) error
	GetAll() ([]T, error)
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
	GetEmptyStruct func() DbType
}

func (uh *CRUDHelper[DbType, ModelType]) Get(id int64) (ModelType, error) {
	dbItem := uh.GetEmptyStruct()
	err := uh.db.Model(dbItem).Where("id=?", id).Select()
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
	dbItem := uh.GetEmptyStruct()
	_, err := uh.db.Model(dbItem).Where("id=?", id).Delete()
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

func (uh *CRUDHelper[DbType, ModelType]) GetAll() ([]ModelType, error) {
	var dbItem []DbType
	var out []ModelType
	err := uh.db.Model(&dbItem).Select()
	if err != nil {
		return out, err
	}

	out = make([]ModelType, 0, len(dbItem))
	for _, item := range dbItem {
		out = append(out, uh.MapModelFromDB(item))
	}

	return out, nil
}
