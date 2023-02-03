package applog

import (
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/DmitryVesenniy/go-rest-framework/framework/authentication"
	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"
	"github.com/DmitryVesenniy/go-rest-framework/framework/typesdb"
)

type ActionType uint

const (
	Login ActionType = iota + 1
	Logout
	DataGet
	DataUpdate
	DataCreate
	DataDelete
)

type AppLog struct {
	models.BaseModel
	serializers.Serializer
	ID        uint                `json:"id" serializer:"read_only" gorm:"primarykey"`
	CreatedAt time.Time           `json:"createdAt" serializer:"read_only;title:Дата и время;"`
	UserID    uint                `json:"userId" serializer:"read_only;" gorm:"column:user_id"`
	Action    ActionType          `json:"action" serializer:"read_only;title:Действие;" gorm:"column:action"`
	Module    string              `json:"module" serializer:"read_only;title:Модуль;" gorm:"column:module"`
	Table     string              `json:"table" serializer:"read_only;title:Таблица;" gorm:"column:table"`
	User      authentication.User `json:"user" serializer:"read_only;title:Пользователь;referenceField:email;" gorm:"foreignkey:UserID;references:ID;OnDelete:SET NULL;"`
	Diff      typesdb.JSONB       `json:"diff" serializer:"read_only;title:Изменения;"  gorm:"type:jsonb;default:'{}';not null;column:diff"`
}

func (AppLog) TableName() string {
	return "log"
}

func (log *AppLog) Create(db *gorm.DB) error {
	tx := db.Omit("User").Create(log)
	return tx.Error
}

func SetAction(opt AppLogOptions, diff DiffInterface) error {
	action := GetAction(opt.Method)

	if action == 0 || opt.User == nil || opt.User.ID == 0 {
		return nil
	}

	diffDict := make(map[string]interface{})
	if diff != nil {
		err := diff.CalcDifference()
		if err != nil {
			return err
		}

		diffDict = diff.ToDict()
	}

	appLogInstance := AppLog{
		Action: action,
		Module: opt.Module,
		Table:  opt.Table,
		UserID: opt.User.ID,
		Diff:   diffDict,
	}

	return appLogInstance.Create(opt.DB)
}

func GetAction(httpMethod string) ActionType {
	var action ActionType

	switch httpMethod {
	case http.MethodPost:
		action = DataCreate
	case http.MethodPut:
		action = DataUpdate
	case http.MethodDelete:
		action = DataDelete
	}

	return action
}
