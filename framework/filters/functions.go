package filters

import (
	"database/sql"
	"time"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"gorm.io/gorm"
)

type FilterFunc func(*gorm.DB, string, string) *gorm.DB

func WhereFilter(param string) FilterFunc {
	return func(queryset *gorm.DB, filterWay, value string) *gorm.DB {
		_where, _val := FilterWay(filterWay)(param, value)
		return queryset.Where(_where, _val)
	}
}

func (f *FilterOption) ValidateQueryParams(value string) bool {
	if f.TypeQueryParam != "" {
		switch f.TypeQueryParam {
		case Number:
			return common.IsDigitsOnly(value)
		case Boolean:
			return value == "1" || value == "true" || value == "0" || value == "false"
		case Date:
			format := "2006-01-02T00:00Z"
			_, err := time.Parse(format, value)
			return err == nil
		}
	}
	return true
}

func CreateManyORParams(params []sql.NamedArg, where string, query *gorm.DB) *gorm.DB {
	if len(params) > 0 {
		query = query.Where(where, params[0])
		for _, param := range params[1:] {
			query = query.Or(query.Where(where, param))
		}
	}
	return query
}
