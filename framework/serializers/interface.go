package serializers

import "github.com/DmitryVesenniy/go-rest-framework/framework/models"

type SerializersInterface interface {
	Validate() bool
	Data() map[string]interface{}
	AddError(field string, err ...string)
	SetError(err string)
	GetModel() models.BaseModelInterface
	GetInstanceModel() models.BaseModelInterface
	Errors() map[string][]string
}

type SerializerModels interface {
	SerializersInterface
	models.BaseModelInterface
}
