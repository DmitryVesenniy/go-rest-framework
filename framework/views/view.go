package views

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/applog"
	"github.com/DmitryVesenniy/go-rest-framework/framework/authentication"
	"github.com/DmitryVesenniy/go-rest-framework/framework/filters"
	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
	"github.com/DmitryVesenniy/go-rest-framework/framework/permissions"
	"github.com/DmitryVesenniy/go-rest-framework/framework/resterrors"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"
)

type ViewType string

const (
	// способ отображения данных
	// на данный момент это либо таблица, либо дерево
	Tree  ViewType = "tree"
	Table ViewType = "table"
	File  ViewType = "file"
	Api   ViewType = "api"
)

type BaseView struct {
	Methods           []string
	Parameters        filters.Parameters
	ViewType          ViewType
	FilteredFields    []filters.FilterQuerySet
	SortAllowFields   []filters.FilterQuerySet
	Serializer        serializers.SerializerModels
	Permissions       []permissions.PermissionInterface
	Module            string
	AppContext        *appctx.AppContext
	PreloadFieldsName []string
	Omit              []string
	AppLog            func(*appctx.AppContext, string, string) error
}

func (b *BaseView) GetID(r *http.Request) (int, error) {
	var vars map[string]string = mux.Vars(r)
	pk := vars["pk"]

	id, err := strconv.Atoi(pk)
	return id, err
}
func (b *BaseView) GetPermissionModule() string {
	return b.Module
}

func (b *BaseView) GetFilteredFields() []filters.FilterQuerySet {
	return b.FilteredFields
}
func (b *BaseView) GetSortAllowFields() []filters.FilterQuerySet {
	return b.SortAllowFields
}

func (b *BaseView) GetSerializer() serializers.SerializerModels {
	return b.Serializer
}

func (b *BaseView) GetViewType() ViewType {
	return Table
}

func (b *BaseView) GetOmit() []string {
	return b.Omit
}

func (b *BaseView) GetPreloadFieldName() []string {
	return b.PreloadFieldsName
}

func (b *BaseView) FilterAndSort(tx *gorm.DB, params filters.Parameters) *gorm.DB {
	for _, filterParameter := range params.FilterOption {
		filterWay, ok := filters.FilterAssociations[filterParameter.FilteringMethod]
		if !ok {
			filterWay = "="
		}

		funcFilter := filters.FilterWay(filterWay)
		val, _ := filterParameter.Value.(string)
		condition, val := funcFilter(filterParameter.CollName, val)
		tx = tx.Where(condition, val)
	}

	if params.SortOption.Field != "" {
		tx = tx.Order(fmt.Sprintf("\"%s\" %s", params.SortOption.CollName, params.SortOption.SortKey))
	}

	return tx
}

func (b *BaseView) GetMethods() []ViewSetMethod {
	return nil
}

func (b *BaseView) GetMiddlewares() []func(http.Handler) http.Handler {
	return nil
}

// HTTP METHODS
func (b *BaseView) Get(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) Post(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) Put(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) LogSetter(appCtx *appctx.AppContext, diff applog.DiffInterface) error {
	model := b.GetSerializer()

	if model != nil {
		tableName := model.TableName()
		module := b.Module

		return applog.SetAction(appCtx, tableName, module, diff)
	}
	return nil
}

// VIEWSET METHODS
func (b *BaseView) List(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) Create(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) Retrieve(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) Update(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

func (b *BaseView) Delete(appCtx *appctx.AppContext) {
	resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
}

// Permissions
func (b *BaseView) CheckPermission(appContext *appctx.AppContext) bool {
	isPermission := true
	if appContext == nil {
		return isPermission
	}

	if len(b.Permissions) > 0 {
		if appContext.User == nil {
			var err error
			appContext.User, err = authentication.GetUserFromContext(appContext.Request.Context().Value(common.ContextUserKey), appContext.DB)
			if err != nil {
				return false
			}
		}
		for _, permission := range b.Permissions {
			if !permission.HasPermission(appContext) {
				isPermission = false
			}
		}
	}
	return isPermission
}
func (b *BaseView) CheckPermissionObject(appContext *appctx.AppContext, model models.BaseModelInterface) bool {
	isPermission := true
	if appContext == nil {
		return isPermission
	}

	if len(b.Permissions) > 0 {
		if appContext.User == nil {
			var err error
			appContext.User, err = authentication.GetUserFromContext(appContext.Request.Context().Value(common.ContextUserKey), appContext.DB)
			if err != nil {
				return false
			}
		}
		for _, permission := range b.Permissions {
			if !permission.HasObjectPermission(appContext, model) {
				isPermission = false
			}
		}
	}
	return isPermission
}
