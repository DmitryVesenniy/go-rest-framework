package autodocs

import (
	"github.com/DmitryVesenniy/go-rest-framework/framework/views"
)

type ReferenceData struct {
	PresentationField string `json:"presentationField"`
}

type ReferenceTable struct {
	TableName string `json:"tableName"`
	PkField   string `json:"pkField"`
}

type SerializerField struct {
	NameField      string          `json:"nameField"`
	Title          string          `json:"title"`
	Type           string          `json:"type"`
	Signatura      interface{}     `json:"signatura"`
	Required       bool            `json:"required"`
	ReadOnly       bool            `json:"readOnly"`
	WriteOnly      bool            `json:"writeOnly"`
	ReferenceData  *ReferenceData  `json:"referenceData"`
	ReferenceTable *ReferenceTable `json:"referenceTable"`
	ProxyField     *ReferenceData  `json:"proxyField"`
}

type ViewData struct {
	ViewType          views.ViewType        `json:"viewType"`
	MethodsAllow      []views.ViewSetMethod `json:"methodsAllow"`
	TableName         string                `json:"tableName"`
	FilterFields      []FilterField         `json:"filterFields"`
	SortAllowFields   []string              `json:"sortAllowFields"`
	SerializersFields []SerializerField     `json:"serializersFields"`
	Description       string                `json:"description"`
	PermissionModule  string                `json:"permissionModule"`
}

type FilterField struct {
	ParamName      string      `json:"paramName"`
	TypeQueryParam interface{} `json:"typeQueryParam"`
	Enum           interface{} `json:"enum"`
	Directory      string      `json:"directory"`
}
