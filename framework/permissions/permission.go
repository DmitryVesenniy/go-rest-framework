package permissions

import (
	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
)

type And struct {
	Permissions []PermissionInterface
}

func (permission *And) HasPermission(appCtx *appctx.AppContext) bool {
	for _, _permission := range permission.Permissions {
		if !_permission.HasPermission(appCtx) {
			return false
		}
	}

	return true
}
func (permission *And) HasObjectPermission(appCtx *appctx.AppContext, model models.BaseModelInterface) bool {
	for _, _permission := range permission.Permissions {
		if !_permission.HasObjectPermission(appCtx, model) {
			return false
		}
	}

	return true
}

type Or struct {
	Permissions []PermissionInterface
}

func (permission *Or) HasPermission(appCtx *appctx.AppContext) bool {
	for _, _permission := range permission.Permissions {
		if _permission.HasPermission(appCtx) {
			return true
		}
	}

	return false
}
func (permission *Or) HasObjectPermission(appCtx *appctx.AppContext, model models.BaseModelInterface) bool {
	for _, _permission := range permission.Permissions {
		if _permission.HasObjectPermission(appCtx, model) {
			return true
		}
	}

	return false
}

type AllowAny struct{}

func (permission *AllowAny) HasPermission(appCtx *appctx.AppContext) bool {
	return true
}
func (permission *AllowAny) HasObjectPermission(appCtx *appctx.AppContext, model models.BaseModelInterface) bool {
	return true
}

type IsAuthenticated struct {
}

func (permission *IsAuthenticated) HasPermission(appCtx *appctx.AppContext) bool {
	return appCtx.User != nil
}
func (permission *IsAuthenticated) HasObjectPermission(appCtx *appctx.AppContext, model models.BaseModelInterface) bool {
	return appCtx.User != nil
}

type IsSuperUser struct {
}

func (permission *IsSuperUser) HasPermission(appCtx *appctx.AppContext) bool {
	return appCtx.User != nil && appCtx.User.IsSuperuser
}
func (permission *IsSuperUser) HasObjectPermission(appCtx *appctx.AppContext, model models.BaseModelInterface) bool {
	return appCtx.User != nil && appCtx.User.IsSuperuser
}
