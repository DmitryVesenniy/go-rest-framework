package appctx

import (
	"context"
	"io"
	"net/http"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"github.com/DmitryVesenniy/go-rest-framework/framework/authentication"

	"gorm.io/gorm"
)

type AppContext struct {
	DB       *gorm.DB
	User     *authentication.User
	Request  *http.Request
	Response http.ResponseWriter
	Ctx      context.Context
	Body     []byte
}

func NewAppContext(w http.ResponseWriter, r *http.Request, db *gorm.DB) *AppContext {
	userRequests, err := authentication.GetUserFromContext(r.Context().Value(common.ContextUserKey), db)
	if err != nil {
		userRequests = nil
	}

	body, _ := io.ReadAll(r.Body)
	r.Body.Close()

	appContext := &AppContext{
		Request:  r,
		Response: w,
		User:     userRequests,
		Body:     body,
		DB:       db,
	}

	return appContext
}
