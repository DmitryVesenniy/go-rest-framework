package framework

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/DmitryVesenniy/go-rest-framework/framework/appctx"
	"github.com/DmitryVesenniy/go-rest-framework/framework/logger"
	"github.com/DmitryVesenniy/go-rest-framework/framework/media"
	"github.com/DmitryVesenniy/go-rest-framework/framework/notifications"
	"github.com/DmitryVesenniy/go-rest-framework/framework/views"
)

// App has router and db instances
type App struct {
	Router        *mux.Router
	DB            *gorm.DB
	MediaService  *media.MediaServise
	LogService    logger.LoggerInterface
	NotivyService notifications.NotificationsInterface
}

// Initialize initializes the app with predefined configuration
func NewApp(db *gorm.DB, notivyService notifications.NotificationsInterface) *App {
	return &App{
		DB:            db,
		NotivyService: notivyService,
		Router:        mux.NewRouter(),
	}
}

func (app *App) wrapperHandlers(hundler func(*appctx.AppContext)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		appCtx := appctx.NewAppContext(w, r, app.DB.Order("id"))
		hundler(appCtx)
	}
}

func (app *App) RegisterGenericView(path string, baseRout *mux.Router, view views.ViewInterface) {
	middlevares := view.GetMiddlewares()
	methods := view.GetMethods()

	for _, middlevare := range middlevares {
		baseRout.Use(middlevare)
	}

	for _, method := range methods {
		switch method {
		case views.GET:
			baseRout.HandleFunc(path, app.wrapperHandlers(view.Get)).Methods(http.MethodGet)
		case views.POST:
			baseRout.HandleFunc(path, app.wrapperHandlers(view.Post)).Methods(http.MethodPost)
		case views.List:
			baseRout.HandleFunc(fmt.Sprintf("%s/", path), app.wrapperHandlers(view.List)).Methods(http.MethodGet)
		case views.Create:
			baseRout.HandleFunc(fmt.Sprintf("%s/", path), app.wrapperHandlers(view.Create)).Methods(http.MethodPost)
		case views.Retrieve:
			baseRout.HandleFunc(fmt.Sprintf("%s/{pk:.+}", path), app.wrapperHandlers(view.Retrieve)).Methods(http.MethodGet)
		case views.Update:
			baseRout.HandleFunc(fmt.Sprintf("%s/{pk:.+}", path), app.wrapperHandlers(view.Update)).Methods(http.MethodPut)
		case views.Delete:
			baseRout.HandleFunc(fmt.Sprintf("%s/{pk:.+}", path), app.wrapperHandlers(view.Delete)).Methods(http.MethodDelete)
		}
	}
}

func (app *App) RegisterGenericMapView(path string, baseRout *mux.Router, viewsAPI map[views.ViewSetMethod]func(*appctx.AppContext)) {
	for method, fn := range viewsAPI {
		switch method {
		case views.List:
			baseRout.HandleFunc(fmt.Sprintf("%s/", path), app.wrapperHandlers(fn)).Methods(http.MethodGet)
		case views.Create:
			baseRout.HandleFunc(fmt.Sprintf("%s/", path), app.wrapperHandlers(fn)).Methods(http.MethodPost)
		case views.Retrieve:
			baseRout.HandleFunc(fmt.Sprintf("%s/{pk:.+}", path), app.wrapperHandlers(fn)).Methods(http.MethodGet)
		case views.Update:
			baseRout.HandleFunc(fmt.Sprintf("%s/{pk:.+}", path), app.wrapperHandlers(fn)).Methods(http.MethodPut)
		case views.Delete:
			baseRout.HandleFunc(fmt.Sprintf("%s/{pk:.+}", path), app.wrapperHandlers(fn)).Methods(http.MethodDelete)
		}
	}
}

func (app *App) RegisterHundler(methods []string, path string, baseRout *mux.Router, hundler func(*appctx.AppContext)) {
	baseRout.HandleFunc(path, app.wrapperHandlers(hundler)).Methods(methods...)
}
