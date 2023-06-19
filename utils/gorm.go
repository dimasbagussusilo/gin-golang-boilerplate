package utils

import "gorm.io/gorm"

func Transaction(db *gorm.DB, f func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func SortBy(orderBy string, query *gorm.DB) *gorm.DB {
	switch orderBy {
	case "newest":
		query.Order("created_at desc")
	case "oldest":
		query.Order("created_at asc")
	default:
		query.Order("created_at desc")
	}
	return query
}
