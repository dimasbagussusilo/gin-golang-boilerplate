package service

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"gorm.io/gorm"
)

type Entity interface {
	TableName() string
}

type Repository struct {
	entity Entity
	db     *gorm.DB
}

type IRepository interface {
	FindOne(entity Entity, id int, dbTransaction *gorm.DB) error
	FindAll(entities any, pageNum int, pageSize int, options func(db *gorm.DB) *gorm.DB, dbTransaction *gorm.DB) (*utils.Pagination, error)
	Create(form Entity, dbTransaction *gorm.DB) (Entity, error)
	Update(id int, form Entity, dbTransaction *gorm.DB) (Entity, error)
	Delete(id int, dbTransaction *gorm.DB) error
	CustomQuery() *gorm.DB
}

func NewRepository(entity Entity, db *gorm.DB) IRepository {
	return &Repository{
		entity: entity,
		db:     db,
	}
}

func (r *Repository) getDB(dbTransaction *gorm.DB) *gorm.DB {
	if dbTransaction != nil {
		return dbTransaction
	}
	return r.db
}

func (r *Repository) FindOne(entity Entity, id int, dbTransaction *gorm.DB) error {
	db := r.getDB(dbTransaction)

	err := db.Model(r.entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindAll(entities any, pageNum int, pageSize int, options func(db *gorm.DB) *gorm.DB, dbTransaction *gorm.DB) (*utils.Pagination, error) {
	db := r.getDB(dbTransaction)

	var count int64
	query := db.Model(r.entity)
	query = options(query)
	query.Count(&count)
	query = query.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	res := query.Find(entities)
	if res.Error != nil {
		fmt.Printf("error, %+v\n", res.Error)
		return nil, res.Error
	}

	pagination := utils.Paginate(count, pageNum, pageSize)
	return pagination, nil
}

func (r *Repository) Create(form Entity, dbTransaction *gorm.DB) (Entity, error) {
	db := r.getDB(dbTransaction)

	result := db.Table(r.entity.TableName()).Select("*").Create(form)
	if result.Error != nil {
		return nil, result.Error
	}

	return form, nil
}

func (r *Repository) Update(id int, form Entity, dbTransaction *gorm.DB) (Entity, error) {
	db := r.getDB(dbTransaction)

	entity := r.entity
	err := db.Model(entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return nil, err
	}

	result := db.Save(form)
	if result.Error != nil {
		return nil, result.Error
	}

	return entity, nil
}

func (r *Repository) Delete(id int, dbTransaction *gorm.DB) error {
	db := r.getDB(dbTransaction)

	entity := r.entity
	err := db.Model(entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return err
	}

	result := db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Repository) CustomQuery() *gorm.DB {
	return r.db
}
