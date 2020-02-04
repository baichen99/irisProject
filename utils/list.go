package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
)

// GetListParameters base parameters needed to define the get list queries
type GetListParameters struct {
	Page    int
	Limit   int
	OrderBy string
	Order   string
}

// GetUserListParameters parameters for GetUserList queries
type GetUserListParameters struct {
	GetListParameters
	Username string
}

// GetProfileListParameters parameters for GetProfileList queries
type GetProfileListParameters struct {
	GetListParameters
	Content string
}

// SearchByColumn 模糊搜索
func SearchByColumn(columnName string, searchString string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchString != "" {
			// add % to prefix and suffix to match search
			// SQL search: SELECT * FROM column_name WHERE query LIKE value
			value := "%" + searchString + "%"
			query := columnName + " LIKE ?"
			return db.Where(query, value)
		}
		return db
	}
}

// FilterByColumn 精确搜索
func FilterByColumn(columnName string, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// SQL select: SELECT * FROM table WHERE columnName = value
		if value != "" {
			query := columnName + " = ?"
			return db.Where(query, value)
		}
		return db

	}
}

// GetListParamsFromContext gets list params from context
func GetListParamsFromContext(c iris.Context, orderName string) (listParmas GetListParameters, err error) {
	listParmas.Page = c.URLParamIntDefault("page", 1)
	listParmas.Limit = c.URLParamIntDefault("limit", 1)
	listParmas.Order = c.URLParamDefault("order", "asc")
	listParmas.OrderBy = c.URLParamDefault("orderBy", orderName)
	return
}
