package models

import (
	"fmt"
	"strings"

	"elogika.vsb.cz/backend/initializers"
	"elogika.vsb.cz/backend/modules/common"
	"gorm.io/gorm"
)

type CommonModel struct{}

func (CommonModel) ApplySorting(query *gorm.DB, sortings []common.SearchRequestSorting) *gorm.DB {
	for _, s := range sortings {
		direction := "ASC"
		if s.Desc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", s.ID, direction))
	}
	return query
}

func (CommonModel) ApplyPagination(query *gorm.DB, pagination *common.SearchRequestPagination) *gorm.DB {
	if pagination != nil {
		pageIndex := pagination.PageIndex
		pageSize := pagination.PageSize

		offset := pageIndex * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	return query
}

func (q CommonModel) ApplyFilters(query *gorm.DB, filters []common.SearchRequestFilter, model any, extra map[string]interface{}) (*gorm.DB, *common.ErrorResponse) {
	for _, filter := range filters {
		column, err := GetModelColumnName(initializers.DB, model, CapitalizeFirstLetter(filter.ID))
		if err != nil {
			// return nil, &common.ErrorResponse{
			// 	Message: "Failed to find filters",
			// 	Details: err.Error(),
			// }
			continue //Todo not the best ????
		}
		value := filter.Value
		// TODO figure out a syntax for defining type of filter (>, <, contains, ...)
		query = query.Where(fmt.Sprintf("%s like ?", column), fmt.Sprintf("%%%v%%", value))
	}
	return query, nil
}

func (CommonModel) GetCount(query *gorm.DB) int64 {
	var totalCount int64 = 0
	query.Count(&totalCount)
	return totalCount
}

func GetModelColumnName(db *gorm.DB, model any, fieldName string) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return "", err
	}

	// Replace only the last "Id" with "ID"
	if len(fieldName) >= 2 && fieldName[len(fieldName)-2:] == "Id" {
		fieldName = fieldName[:len(fieldName)-2] + "ID"
	}

	field := stmt.Schema.LookUpField(fieldName)
	if field == nil {
		return "", fmt.Errorf("Field %s not found in schema", fieldName)
	}
	return field.DBName, nil
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
