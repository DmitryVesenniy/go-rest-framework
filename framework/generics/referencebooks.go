package generics

import (
	"github.com/gorilla/mux"

	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/autodocs"
	"github.com/DmitryVesenniy/go-rest-framework/framework/permissions"
	"github.com/DmitryVesenniy/go-rest-framework/framework/response"
	"github.com/DmitryVesenniy/go-rest-framework/framework/resterrors"
	"github.com/DmitryVesenniy/go-rest-framework/framework/views"
)

type ReferenceBookView struct {
	views.BaseView
	TablesAPIViews        map[string]views.ViewInterface
	TablesGenericAPIViews map[string]map[views.ViewSetMethod]func(*appctx.AppContext)
	Docs                  []autodocs.ViewData
}

func (*ReferenceBookView) GetMethods() []views.ViewSetMethod {
	return []views.ViewSetMethod{
		views.List,
		views.Create,
		views.Retrieve,
		views.Update,
		views.Delete,
		views.GET,
	}
}

func (rb *ReferenceBookView) List(appCtx *appctx.AppContext) {
	if !rb.CheckPermission(appCtx) {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}

	response.JsonRespond(appCtx.Response, map[string]interface{}{"data": rb.Docs})
}

func (rb *ReferenceBookView) Get(appCtx *appctx.AppContext) {
	if !rb.CheckPermission(appCtx) {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	vars := mux.Vars(appCtx.Request)
	table := vars["table"]

	view, ok := rb.TablesAPIViews[table]
	if !ok || view == nil {
		hundlerFns, ok := rb.TablesGenericAPIViews[table]
		if !ok || hundlerFns == nil {
			resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
			return
		}
		fn, isOk := hundlerFns[views.List]
		if !isOk || fn == nil {
			resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
			return
		}

		fn(appCtx)
		return
	}

	view.List(appCtx)
}

func (rb *ReferenceBookView) Create(appCtx *appctx.AppContext) {
	if !rb.CheckPermission(appCtx) {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	vars := mux.Vars(appCtx.Request)
	table := vars["table"]

	view, ok := rb.TablesAPIViews[table]
	if !ok || view == nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	view.Create(appCtx)
}

func (rb *ReferenceBookView) Retrieve(appCtx *appctx.AppContext) {
	if !rb.CheckPermission(appCtx) {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	vars := mux.Vars(appCtx.Request)
	table := vars["table"]

	view, ok := rb.TablesAPIViews[table]

	if !ok || view == nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	view.Retrieve(appCtx)
}

func (rb *ReferenceBookView) Update(appCtx *appctx.AppContext) {
	if !rb.CheckPermission(appCtx) {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	vars := mux.Vars(appCtx.Request)
	table := vars["table"]

	view, ok := rb.TablesAPIViews[table]

	if !ok || view == nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	view.Update(appCtx)
}

func (rb *ReferenceBookView) Delete(appCtx *appctx.AppContext) {
	if !rb.CheckPermission(appCtx) {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.PermissionDeniedErr)
		return
	}
	vars := mux.Vars(appCtx.Request)
	table := vars["table"]

	view, ok := rb.TablesAPIViews[table]

	if !ok || view == nil {
		resterrors.RestErrorResponce(appCtx.Response, resterrors.NotFoundErr)
		return
	}

	view.Delete(appCtx)
}

func NewReferenceBookAPIView(viewList []views.ViewInterface, permissions []permissions.PermissionInterface) *ReferenceBookView {
	apiDocs := autodocs.GenerateJSON(viewList)

	rb := &ReferenceBookView{
		BaseView: views.BaseView{
			Permissions: permissions,
		},
		TablesAPIViews: make(map[string]views.ViewInterface),
		Docs:           apiDocs,
	}

	for _, view := range viewList {
		serializer := view.GetSerializer()
		if serializer != nil {
			if serializer.TableName() != "" {
				rb.TablesAPIViews[serializer.TableName()] = view
			}
		}
	}
	return rb
}
