package models

import "gorm.io/gorm"

type BaseModelInterface interface {
	Create(*gorm.DB) error
	Update(*gorm.DB) error
	Delete(*gorm.DB) error
	Description() string
	TableName() string
	GetPk() uint
}
