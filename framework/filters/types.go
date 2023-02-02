package filters

import "gorm.io/gorm"

type SortKey string
type FilteringMethod string
type TypeQueryParam string

const (
	Number  TypeQueryParam = "number"
	Boolean TypeQueryParam = "bool"
	String  TypeQueryParam = "string"
	Date    TypeQueryParam = "date"
)

type QueryParam struct {
	TypeQueryParam TypeQueryParam `json:"typeQueryParam"`
	Enum           []string       `json:"enum"`
	Directory      string         `json:"directory"`
}

type FilterAndSortInterface interface {
	Filter(*gorm.DB) *gorm.DB
	Sort(*gorm.DB) *gorm.DB
}

type FilterViewInterface interface {
	GetFilteredFields() []FilterQuerySet
	GetSortAllowFields() []FilterQuerySet
}

type FilterQuerySet struct {
	ParamName  string
	FieldName  string
	QueryParam QueryParam
	FilterFunc func(*gorm.DB, interface{}) *gorm.DB
}

type SortOption struct {
	Field    string
	CollName string
	SortKey  SortKey
}

type FilterOption struct {
	Field           string
	Value           interface{}
	CollName        string
	FilteringMethod FilteringMethod
	FilterFunc      func(*gorm.DB, interface{}) *gorm.DB
	TypeQueryParam  TypeQueryParam
}

type Parameters struct {
	FilterOption []FilterOption
	SortOption   SortOption
}

type FilterParameter struct {
	Value       string
	Modificator string
	Param       string
}
