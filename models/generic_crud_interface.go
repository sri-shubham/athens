package models

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sri-shubham/athens/util"
)

// Generic CRUD operations interface
type CRUD[T any] interface {
	GetNewEmptyStruct() T
	Get(ctx context.Context, id int64) (T, error)
	Create(context.Context, T) (int64, error)
	Update(context.Context, T) (int64, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]T, error)
	GetBulk(ctx context.Context, ids []int64) ([]T, error)
}

type IdGetter interface {
	GetID() int64
}

// We are adding generic Implementation so we can
// to reduce code duplication
type CRUDHelper[DbType IdGetter, ModelType IdGetter] struct {
	db             *pg.DB
	updateQueue    util.UpdateQueue
	MapModelToDB   func(ModelType) DbType
	MapModelFromDB func(DbType) ModelType
	GetEmptyStruct func() DbType
}

func (uh *CRUDHelper[DbType, ModelType]) Get(ctx context.Context, id int64) (ModelType, error) {
	dbItem := uh.GetEmptyStruct()
	err := uh.db.Model(dbItem).Context(ctx).Where("id=?", id).Select()
	if err != nil {
		var v ModelType
		return v, err
	}
	return uh.MapModelFromDB(dbItem), nil
}

func (uh *CRUDHelper[DbType, ModelType]) Create(ctx context.Context, u ModelType) (int64, error) {
	dbItem := uh.MapModelToDB(u)
	_, err := uh.db.Model(dbItem).Context(ctx).Insert()
	if err != nil {
		return 0, err
	}

	uh.updateQueue.Enqueue(ctx, &util.Item{
		Action: util.ActionCreate,
		ID:     dbItem.GetID(),
		Value:  uh.MapModelFromDB(dbItem),
	})

	return dbItem.GetID(), nil
}

func (uh *CRUDHelper[DbType, ModelType]) Delete(ctx context.Context, id int64) error {
	dbItem := uh.GetEmptyStruct()
	err := uh.db.Model(dbItem).Context(ctx).Where("id=?", id).Select()
	if err != nil {
		return err
	}

	dbItem2 := uh.GetEmptyStruct()
	_, err = uh.db.Model(dbItem2).Context(ctx).Where("id=?", id).Delete()
	if err != nil {
		return err
	}

	uh.updateQueue.Enqueue(ctx, &util.Item{
		Action: util.ActionDelete,
		ID:     dbItem.GetID(),
		Value:  uh.MapModelFromDB(dbItem),
	})

	return nil
}

func (uh *CRUDHelper[DbType, ModelType]) Update(ctx context.Context, in ModelType) (int64, error) {
	oldDbItem := uh.GetEmptyStruct()
	err := uh.db.Model(oldDbItem).Context(ctx).Where("id=?", in.GetID()).Select()
	if err != nil {
		return 0, err
	}

	dbItem := uh.MapModelToDB(in)
	_, err = uh.db.Model(dbItem).Context(ctx).WherePK().Update()
	if err != nil {
		return 0, err
	}

	uh.updateQueue.Enqueue(ctx, &util.Item{
		Action:   util.ActionUpdate,
		ID:       dbItem.GetID(),
		Value:    uh.MapModelFromDB(dbItem),
		OldValue: uh.MapModelFromDB(oldDbItem),
	})

	return dbItem.GetID(), nil
}

func (uh *CRUDHelper[DbType, ModelType]) GetAll(ctx context.Context) ([]ModelType, error) {
	var dbItem []DbType
	var out []ModelType
	err := uh.db.Model(&dbItem).Context(ctx).Select()
	if err != nil {
		return out, err
	}

	out = make([]ModelType, 0, len(dbItem))
	for _, item := range dbItem {
		out = append(out, uh.MapModelFromDB(item))
	}

	return out, nil
}

func (uh *CRUDHelper[DbType, ModelType]) GetBulk(ctx context.Context, ids []int64) ([]ModelType, error) {
	var out []ModelType
	if len(ids) == 0 {
		return out, nil
	}

	var dbItem []DbType

	err := uh.db.Model(&dbItem).Context(ctx).Where("id in (?)", pg.In(ids)).Select()
	if err != nil {
		return out, err
	}

	out = make([]ModelType, 0, len(dbItem))
	for _, item := range dbItem {
		out = append(out, uh.MapModelFromDB(item))
	}

	return out, nil
}
