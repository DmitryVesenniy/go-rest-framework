package views

import (
	"net/http"

	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/applog"
	"github.com/DmitryVesenniy/go-rest-framework/framework/filters"
	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"

	"gorm.io/gorm"
)

type ViewInterface interface {
	GetMethods() []ViewSetMethod
	GetMiddlewares() []func(http.Handler) http.Handler
	FilterAndSort(*gorm.DB, filters.Parameters) *gorm.DB
	GetViewType() ViewType
	GetID(r *http.Request) (int, error)
	GetPermissionModule() string

	Get(*appctx.AppContext)
	Post(*appctx.AppContext)
	Put(*appctx.AppContext)

	List(*appctx.AppContext)
	Create(*appctx.AppContext)
	Retrieve(*appctx.AppContext)
	Update(*appctx.AppContext)
	Delete(*appctx.AppContext)

	GetSerializer() serializers.SerializerModels
	CheckPermission(*appctx.AppContext) bool
	CheckPermissionObject(*appctx.AppContext, models.BaseModelInterface) bool
	GetPreloadFieldName() []string
	LogSetter(*appctx.AppContext, applog.DiffInterface) error

	GetFilteredFields() []filters.FilterQuerySet
	GetSortAllowFields() []filters.FilterQuerySet
}

type ViewGenericAPIInterface interface {
	ViewInterface
	GetQuerySet(*gorm.DB, *appctx.AppContext) *gorm.DB
	GetSerializerModels() serializers.SerializerModels
	GetSliceSerializerModels() interface{}
	GetOmit() []string
}
