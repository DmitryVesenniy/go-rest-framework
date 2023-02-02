package generics

import (
	"gorm.io/gorm"

	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/filters"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"
	"github.com/DmitryVesenniy/go-rest-framework/framework/views"
)

type GenericView struct {
	views.BaseView
	View views.ViewGenericAPIInterface
}

func (v *GenericView) GetMethods() []views.ViewSetMethod {
	methods := v.View.GetMethods()
	if len(methods) > 0 {
		return methods
	}
	return []views.ViewSetMethod{
		views.List,
		views.Create,
		views.Retrieve,
		views.Update,
		views.Delete,
		views.GET,
	}
}

func (*GenericView) GetViewType() views.ViewType {
	return views.Table
}

func (v *GenericView) GetSerializerModels() serializers.SerializerModels {
	return v.View.GetSerializerModels()
}
func (v *GenericView) GetSliceSerializerModels() interface{} {
	return v.View.GetSliceSerializerModels()
}

func (v *GenericView) GetQuerySet(db *gorm.DB, appCtx *appctx.AppContext) *gorm.DB {
	return v.View.GetQuerySet(db, appCtx)
}

func (v *GenericView) GetPermissionModule() string {
	return v.View.GetPermissionModule()
}

func (v *GenericView) GetFilteredFields() []filters.FilterQuerySet {
	return v.View.GetFilteredFields()
}
func (v *GenericView) GetSortAllowFields() []filters.FilterQuerySet {
	return v.View.GetSortAllowFields()
}

func (v *GenericView) GetSerializer() serializers.SerializerModels {
	return v.View.GetSerializer()
}

func (v *GenericView) List(appCtx *appctx.AppContext) {
	data := v.View.GetSliceSerializerModels()
	views.GenericList(v.View, v.View.GetQuerySet(appCtx.DB, appCtx), data, appCtx)
}

func (v *GenericView) Create(appCtx *appctx.AppContext) {
	data := v.View.GetSerializerModels()
	views.GenericCreate(v.View, v.View.GetQuerySet(appCtx.DB, appCtx), data, appCtx, nil)
}

func (v *GenericView) Retrieve(appCtx *appctx.AppContext) {
	data := v.View.GetSerializerModels()
	views.GenericRetrieve(v.View, v.View.GetQuerySet(appCtx.DB, appCtx), data, appCtx, nil)
}

func (v *GenericView) Update(appCtx *appctx.AppContext) {
	data := v.View.GetSerializerModels()
	views.GenericUpdate(v.View, v.View.GetQuerySet(appCtx.DB, appCtx), data, appCtx, nil)
}

func (v *GenericView) Delete(appCtx *appctx.AppContext) {
	data := v.View.GetSerializerModels()
	views.GenericDelete(v, v.View.GetQuerySet(appCtx.DB, appCtx), data, appCtx, nil)
}
