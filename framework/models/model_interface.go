package models

import "gorm.io/gorm"

type QuerySet interface {
	Get(model BaseModelInterface) error
	Filter(model BaseModelInterface) error
}

type BaseModelInterface interface {
	Create(*gorm.DB) error
	Update(*gorm.DB) error
	Delete(*gorm.DB) error
	Description() string
	TableName() string
	GetPk() interface{}
	SetPK(pk interface{}) error
	PrimaryKey(v interface{}) map[string]interface{}
}
