package filters

import (
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

const (
	REQUES_PARAM_SORT_KEY  = "sort"
	SEPARATOR_METHOD_VALUE = "_"

	ASC  SortKey = "asc"
	DESC SortKey = "desc"

	GT   FilteringMethod = "gt"
	GTE  FilteringMethod = "gte"
	LT   FilteringMethod = "lt"
	LTE  FilteringMethod = "lte"
	EQ   FilteringMethod = "eq"
	IN   FilteringMethod = "in"
	LIKE FilteringMethod = "like"
)

var (
	FilterAssociations map[FilteringMethod]string = map[FilteringMethod]string{
		GT:   ">",
		GTE:  ">=",
		LT:   "<",
		LTE:  "<=",
		EQ:   "=",
		IN:   "IN",
		LIKE: "LIKE",
	}
)

type FilterAndSortBackend struct {
	FilterParams []FilterOption
	SortParams   SortOption
}

func NewFilter(r *http.Request, view FilterViewInterface) *FilterAndSortBackend {
	f := &FilterAndSortBackend{}
	params := FilterFromRequest(r, view.GetFilteredFields(), view.GetSortAllowFields())
	f.FilterParams = params.FilterOption
	f.SortParams = params.SortOption
	return f
}

func (f *FilterAndSortBackend) Filter(query *gorm.DB) *gorm.DB {
	for _, filterParameter := range f.FilterParams {
		val, _ := filterParameter.Value.(string)
		if filterParameter.FilterFunc != nil {
			query = filterParameter.FilterFunc(query, val)
		} else {
			filterWay, ok := FilterAssociations[filterParameter.FilteringMethod]
			if !ok {
				filterWay = "="
			}

			funcFilter := FilterWay(filterWay)

			condition, val := funcFilter(filterParameter.CollName, val)
			query = query.Where(condition, val)
		}
	}

	return query
}
func (f *FilterAndSortBackend) Sort(query *gorm.DB) *gorm.DB {
	if f.SortParams.Field != "" {
		query = query.Order(fmt.Sprintf("\"%s\" %s", f.SortParams.CollName, f.SortParams.SortKey))
	}
	return query
}

func FilterFromRequest(r *http.Request, filterFields []FilterQuerySet, sortAllowFields []FilterQuerySet) Parameters {
	filterParameters := Parameters{}
	allowFilterParams := make(map[string]FilterQuerySet)

	for _, filterField := range filterFields {
		allowFilterParams[filterField.ParamName] = filterField
	}

	for _param, _value := range r.URL.Query() {

		splited := strings.Split(_param, SEPARATOR_METHOD_VALUE)
		modificator := ""
		if len(splited) > 1 {
			modificator = splited[len(splited)-1]
			_param = strings.Join(splited[:len(splited)-1], SEPARATOR_METHOD_VALUE)
		}

		itemfilter, ok := allowFilterParams[_param]

		if !ok {
			continue
		}

		filterOpt := FilterOption{
			Field:      itemfilter.ParamName,
			CollName:   itemfilter.FieldName,
			FilterFunc: itemfilter.FilterFunc,
			Value:      _value[0],
		}

		if modificator == "" {
			filterOpt.FilteringMethod = EQ
		} else {
			filterOpt.FilteringMethod = FilteringMethod(modificator)
		}

		filterParameters.FilterOption = append(filterParameters.FilterOption, filterOpt)
	}

	sortValueParam := r.URL.Query().Get(REQUES_PARAM_SORT_KEY)

	if sortValueParam != "" {
		sortSplited := strings.Split(sortValueParam, SEPARATOR_METHOD_VALUE)
		if len(sortSplited) == 1 {
			filterParameters.SortOption.SortKey = ASC
			filterParameters.SortOption.Field = sortValueParam
		} else {
			filterParameters.SortOption.SortKey = SortKey(sortSplited[0])
			filterParameters.SortOption.Field = strings.Join(sortSplited[1:], SEPARATOR_METHOD_VALUE)
		}
	}

	search := false
	for _, sortField := range sortAllowFields {
		if sortField.ParamName == filterParameters.SortOption.Field {
			search = true
			filterParameters.SortOption.CollName = sortField.FieldName
			break
		}
	}

	if !search {
		filterParameters.SortOption.Field = ""
	}

	if filterParameters.SortOption.SortKey != ASC && filterParameters.SortOption.SortKey != DESC {
		filterParameters.SortOption.SortKey = ASC
	}

	return filterParameters
}

func ConvertFromLikeSearch(word string) string {
	return fmt.Sprintf("%%%s%%", strings.ToLower(word))
}

func FilterWay(filterWay string) func(string, string) (string, string) {
	if filterWay == string(LIKE) {
		return func(collName string, val string) (string, string) {
			return fmt.Sprintf("LOWER(\"%s\") LIKE ?", collName), ConvertFromLikeSearch(val)
		}
	}
	return func(collName string, val string) (string, string) {
		return fmt.Sprintf("\"%s\" %s ?", collName, filterWay), val
	}
}
