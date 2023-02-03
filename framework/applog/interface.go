package applog

import (
	"github.com/DmitryVesenniy/go-rest-framework/framework/authentication"
	"gorm.io/gorm"
)

type DiffInterface interface {
	ToDict() map[string]interface{}
	CalcDifference() error
}

type AppLogDiffInterface interface {
	SetAction(AppLogOptions, DiffInterface) error
	GetAction(string) ActionType
}

type AppLogOptions struct {
	User   *authentication.User
	DB     *gorm.DB
	Table  string
	Module string
	Method string
}
