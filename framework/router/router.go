package router

import (
	"net/http"

	"github.com/DmitryVesenniy/go-rest-framework/framework"
	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
)

func WrapperHandlers(hundler func(*appctx.AppContext), app *framework.App) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := framework.NewAppContext(w, r, app)
		hundler(appCtx)
	}
}
