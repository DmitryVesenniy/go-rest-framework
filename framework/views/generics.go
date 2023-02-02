package views

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/DmitryVesenniy/go-rest-framework/difference"
	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/filters"
	"github.com/DmitryVesenniy/go-rest-framework/framework/pagination"
	"github.com/DmitryVesenniy/go-rest-framework/framework/response"
	"github.com/DmitryVesenniy/go-rest-framework/framework/resterrors"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"

	"gorm.io/gorm"
)

type BringMind struct {
	Before          func(serializers.SerializerModels, *appctx.AppContext) error
	MiddlewareModel func(serializers.SerializerModels, *appctx.AppContext) error
	After           func(serializers.SerializerModels, serializers.SerializerModels, *appctx.AppContext) error
}

func GenericList(view ViewInterface, queryset *gorm.DB, data interface{}, appCtx *appctx.AppContext) {
	isPermission := view.CheckPermission(appCtx)
	if !isPermission {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	filterInstance := filters.NewFilter(appCtx.Request, view)
	queryset = filterInstance.Filter(queryset)
	queryset = filterInstance.Sort(queryset)
	paginator := pagination.PaginatorFromRequest(appCtx.Request, queryset)
	queryset = queryset.Scopes(paginator.Paginate())
	if queryset.Error != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: queryset.Error})
		return
	}

	preloadFields := view.GetPreloadFieldName()
	if len(preloadFields) > 0 {
		for _, preloadName := range preloadFields {
			queryset = queryset.Preload(preloadName)
		}
	}

	err := queryset.Find(data).Error
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: err})
		return
	}

	response.ResponseList(appCtx.Response, data, paginator)
}

func GenericCreate(view ViewInterface, queryset *gorm.DB, data serializers.SerializerModels, appCtx *appctx.AppContext, middle *BringMind) {
	isPermission := view.CheckPermission(appCtx)
	if !isPermission {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	// err := json.NewDecoder(appCtx.Request.Body).Decode(data)
	err := json.Unmarshal(appCtx.Body, data)
	if err != nil {
		restErr := &resterrors.BaseError{Err: err.Error()}
		resterrors.RestErrorResponce(appCtx.Response, restErr)
		return
	}

	if middle != nil && middle.Before != nil {
		err := middle.Before(data, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	if !data.Validate() {
		response.JsonErrorResponce(appCtx.Response, data.Errors(), http.StatusUnprocessableEntity)
		return
	}

	err = data.Create(appCtx.DB)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: err})
		return
	}

	if middle != nil && middle.After != nil {
		err := middle.After(data, nil, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	preloadFields := view.GetPreloadFieldName()
	if len(preloadFields) > 0 {
		for _, preloadName := range preloadFields {
			queryset = queryset.Preload(preloadName)
		}
	}

	err = queryset.First(data).Error

	if err != nil {
		restErr := &resterrors.BaseError{Err: err.Error()}
		resterrors.RestErrorResponce(appCtx.Response, restErr)
		return
	}

	view.LogSetter(appCtx, nil)

	response.Response(appCtx.Response, data)
}

func GenericRetrieve(view ViewInterface, queryset *gorm.DB, data serializers.SerializerModels, appCtx *appctx.AppContext, middle *BringMind) {
	isPermission := view.CheckPermission(appCtx)
	if !isPermission {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	pk, err := view.GetID(appCtx.Request)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	preloadFields := view.GetPreloadFieldName()
	if len(preloadFields) > 0 {
		for _, preloadName := range preloadFields {
			queryset = queryset.Preload(preloadName)
		}
	}

	queryset = queryset.First(data, "id = ?", pk)

	isPermissionObject := view.CheckPermissionObject(appCtx, data)
	if !isPermissionObject {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}

	if queryset.Error != nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	if middle != nil && middle.After != nil {
		err := middle.After(data, nil, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	response.Response(appCtx.Response, data)
}

func GenericUpdate(view ViewInterface, queryset *gorm.DB, data serializers.SerializerModels, appCtx *appctx.AppContext, middle *BringMind) {
	isPermission := view.CheckPermission(appCtx)
	if !isPermission {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	pk, err := view.GetID(appCtx.Request)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	// err = json.NewDecoder(appCtx.Request.Body).Decode(data)
	err = json.Unmarshal(appCtx.Body, data)
	if err != nil {
		restErr := &resterrors.BaseError{Err: err.Error()}
		resterrors.RestErrorResponce(appCtx.Response, restErr)
		return
	}

	valueType := reflect.TypeOf(data)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	dataDB := reflect.New(valueType).Interface().(serializers.SerializerModels)
	err = queryset.First(dataDB, "id = ?", pk).Error
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}
	isPermissionObject := view.CheckPermissionObject(appCtx, dataDB)
	if !isPermissionObject {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}

	if middle != nil && middle.MiddlewareModel != nil {
		err := middle.MiddlewareModel(dataDB, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	err = data.SetPK(pk)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: err})
		return
	}

	if middle != nil && middle.Before != nil {
		err := middle.Before(data, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	diff := difference.NewDifference(pk, dataDB, data)
	errCalcDiff := diff.CalcDifference()

	err = data.Update(appCtx.DB)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: err})
		return
	}

	preloadFields := view.GetPreloadFieldName()
	if len(preloadFields) > 0 {
		for _, preloadName := range preloadFields {
			queryset = queryset.Preload(preloadName)
		}
	}

	dataResponse := reflect.New(valueType).Interface().(serializers.SerializerModels)
	err = queryset.First(dataResponse, "id = ?", pk).Error
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: err})
		return
	}

	if middle != nil && middle.After != nil {
		err := middle.After(dataResponse, dataDB, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	if errCalcDiff == nil {
		view.LogSetter(appCtx, diff)
	}

	response.Response(appCtx.Response, dataResponse)
}

func GenericDelete(view ViewInterface, queryset *gorm.DB, data serializers.SerializerModels, appCtx *appctx.AppContext, middle *BringMind) {
	isPermission := view.CheckPermission(appCtx)
	if !isPermission {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	pk, err := view.GetID(appCtx.Request)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	queryset = queryset.First(data, "id = ?", pk)

	if queryset.Error != nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	isPermissionObject := view.CheckPermissionObject(appCtx, data)
	if !isPermissionObject {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}

	if middle != nil && middle.Before != nil {
		err := middle.Before(data, appCtx)
		if err != nil {
			restErr := &resterrors.BaseError{Err: err.Error()}
			resterrors.RestErrorResponce(appCtx.Response, restErr)
			return
		}
	}

	err = data.Delete(appCtx.DB)
	if err != nil {
		resterrors.RestErrorResponce(appCtx.Response, &resterrors.ModelInstanceError{Err: err})
		return
	}

	view.LogSetter(appCtx, nil)

	response.Response(appCtx.Response, nil)
}

func GenericAPIView(viewAPI ViewGenericAPIInterface) map[ViewSetMethod]func(*appctx.AppContext) {
	methods := viewAPI.GetMethods()
	viewAPIMethods := make(map[ViewSetMethod]func(*appctx.AppContext))

	for _, method := range methods {
		switch method {
		case List:
			viewAPIMethods[List] = func(appCtx *appctx.AppContext) {
				querySet := viewAPI.GetQuerySet(appCtx.DB, appCtx)
				data := viewAPI.GetSliceSerializerModels()
				GenericList(viewAPI, querySet, &data, appCtx)
			}
		case Create:
			viewAPIMethods[Create] = func(appCtx *appctx.AppContext) {
				querySet := viewAPI.GetQuerySet(appCtx.DB, appCtx)
				data := viewAPI.GetSerializerModels()
				GenericCreate(viewAPI, querySet, data, appCtx, nil)
			}
		case Retrieve:
			viewAPIMethods[Retrieve] = func(appCtx *appctx.AppContext) {
				querySet := viewAPI.GetQuerySet(appCtx.DB, appCtx)
				data := viewAPI.GetSerializer()
				GenericRetrieve(viewAPI, querySet, data, appCtx, nil)
			}
		case Update:
			viewAPIMethods[Update] = func(appCtx *appctx.AppContext) {
				querySet := viewAPI.GetQuerySet(appCtx.DB, appCtx)
				data := viewAPI.GetSerializer()
				GenericUpdate(viewAPI, querySet, data, appCtx, nil)
			}
		case Delete:
			viewAPIMethods[Delete] = func(appCtx *appctx.AppContext) {
				querySet := viewAPI.GetQuerySet(appCtx.DB, appCtx)
				data := viewAPI.GetSerializer()
				GenericDelete(viewAPI, querySet, data, appCtx, nil)
			}
		}
	}

	return viewAPIMethods
}
