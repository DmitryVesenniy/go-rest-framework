package permissions

import (
	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
)

type PermissionInterface interface {
	HasPermission(*appctx.AppContext) bool
	HasObjectPermission(*appctx.AppContext, models.BaseModelInterface) bool
}
