package autodocs

import (
	"github.com/DmitryVesenniy/go-rest-framework/framework/views"
)

func GenerateJSON(viewList []views.ViewInterface) []ViewData {
	docs := make([]ViewData, 0, len(viewList))

	for _, view := range viewList {
		tableName := ""
		serializer := view.GetSerializer()
		if serializer != nil {
			tableName = serializer.TableName()
		}

		filtersParams := view.GetFilteredFields()
		filtersFileds := make([]FilterField, 0, len(filtersParams))
		for _, filterParam := range filtersParams {
			filterField := FilterField{
				ParamName:      filterParam.ParamName,
				Directory:      filterParam.QueryParam.Directory,
				TypeQueryParam: filterParam.QueryParam.TypeQueryParam,
				Enum:           filterParam.QueryParam.Enum,
			}
			filtersFileds = append(filtersFileds, filterField)
		}

		sortParams := view.GetSortAllowFields()
		sortFields := make([]string, 0, len(sortParams))
		for _, sortParam := range sortParams {
			sortFields = append(sortFields, sortParam.ParamName)
		}

		doc := ViewData{
			ViewType:         view.GetViewType(),
			MethodsAllow:     view.GetMethods(),
			TableName:        tableName,
			FilterFields:     filtersFileds,
			SortAllowFields:  sortFields,
			PermissionModule: view.GetPermissionModule(),
		}

		if serializer != nil {
			doc.Description = serializer.Description()
		}
		doc.SerializersFields = GetSerializerFields(serializer, 0)
		docs = append(docs, doc)
	}

	return docs
}
