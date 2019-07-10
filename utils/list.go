package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
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

// SearchByColumn 模糊搜索
func SearchByColumn(column string, searchString string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchString != "" {
			queryString := column + " LIKE ?"
			value := "% " + searchString + " %"
			db.Where(queryString, value)
		}
		return db
	}
}

// FilterByColumn 精确搜索
func FilterByColumn(columnName string, value string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// SQL select: SELECT * FROM column_name WHERE query = value
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
	listParmas.Limit = c.URLParamIntDefault("limit", 0) // 0 means no limit
	listParmas.Order = c.URLParamDefault("order", "asc")
	listParmas.OrderBy = c.URLParamDefault("orderBy", orderName)
	return
}
