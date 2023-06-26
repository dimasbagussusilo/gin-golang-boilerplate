package service

import (
	"fmt"
	"github.com/dimasbagussusilo/gin-golang-boilerplate/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Entity interface {
	TableName() string
}

type Repository struct {
	entity Entity
	db     *gorm.DB
}

type IRepository interface {
	FindOne(ctx *gin.Context, entity Entity, id int, dbTransaction *gorm.DB) error
	FindAll(ctx *gin.Context, entities any, options *OptionQuery, dbTransaction *gorm.DB) (*utils.Pagination, error)
	Create(ctx *gin.Context, form Entity, dbTransaction *gorm.DB) (Entity, error)
	Update(ctx *gin.Context, id int, form Entity, dbTransaction *gorm.DB) (Entity, error)
	Delete(ctx *gin.Context, id int, dbTransaction *gorm.DB) error
	CustomQuery(ctx *gin.Context) *gorm.DB
}

type OptionQuery struct {
	Page   int
	Limit  int
	Filter map[string]any
	Search []string
	Order  string
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

func (r *Repository) FindOne(ctx *gin.Context, entity Entity, id int, dbTransaction *gorm.DB) error {
	db := r.getDB(dbTransaction)

	err := db.Model(r.entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindAll(ctx *gin.Context, entities any, options *OptionQuery, dbTransaction *gorm.DB) (*utils.Pagination, error) {
	page := 1
	limit := 10

	db := r.getDB(dbTransaction)

	var count int64
	query := db.Model(r.entity)

	if options != nil {
		if options.Page != 0 {
			page = options.Page
		}
		if options.Limit != 0 {
			limit = options.Limit
		}

		if len(options.Filter) != 0 {
			if customFields, ok := options.Filter["custom_fields"]; ok {
				for i := range customFields.([]clause.Expr) {
					query = query.Where(customFields.([]clause.Expr)[i])
				}
				delete(options.Filter, "custom_fields")
			}
			query = query.Where(options.Filter)
		}

		searchParam := ctx.Query("search")
		if searchParam != "" {
			if options.Search != nil {
				var searchQueries string
				for i := range options.Search {
					if strings.ContainsAny(searchParam, ";") {
						searchParam = strings.Join(strings.Split(searchParam, ";"), "")
					}

					if i == 0 {
						searchQueries += fmt.Sprintf("LOWER(%v) LIKE %v", options.Search[i], "'%"+strings.ToLower(searchParam)+"%'")
					} else {
						searchQueries += fmt.Sprintf(" OR LOWER(%v) LIKE %v", options.Search[i], "'%"+strings.ToLower(searchParam)+"%'")
					}
				}

				query = query.Where(searchQueries)
			}
		}
	}

	query.Count(&count)

	query = query.Limit(limit).Offset((page - 1) * limit)

	if options != nil && options.Order != "" {
		query = utils.SortBy(options.Order, query)
	}

	res := query.Find(entities)
	if res.Error != nil {
		fmt.Printf("error, %+v\n", res.Error)
		return nil, res.Error
	}

	pagination := utils.Paginate(count, page, limit)
	return pagination, nil
}

func (r *Repository) Create(ctx *gin.Context, form Entity, dbTransaction *gorm.DB) (Entity, error) {
	db := r.getDB(dbTransaction)

	result := db.Table(r.entity.TableName()).Select("*").Create(form)
	if result.Error != nil {
		return nil, result.Error
	}

	return form, nil
}

func (r *Repository) Update(ctx *gin.Context, id int, form Entity, dbTransaction *gorm.DB) (Entity, error) {
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

func (r *Repository) Delete(ctx *gin.Context, id int, dbTransaction *gorm.DB) error {
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

func (r *Repository) CustomQuery(ctx *gin.Context) *gorm.DB {
	return r.db
}
