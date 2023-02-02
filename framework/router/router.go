package router

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
)

func WrapperHandlers(hundler func(*appctx.AppContext), db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := appctx.NewAppContext(w, r, db.Order("id"))
		hundler(appCtx)
	}
}
