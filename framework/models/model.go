package models

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	Desc string `json:"-" gorm:"-"`
}

func (bm *BaseModel) Create(db *gorm.DB) error {
	return nil
}
func (bm *BaseModel) Update(db *gorm.DB) error {
	return nil
}
func (bm *BaseModel) Delete(db *gorm.DB) error {
	return nil
}
func (bm *BaseModel) Description() string {
	return bm.Desc
}
func (bm *BaseModel) TableName() string {
	return ""
}

func (bm *BaseModel) GetPk() uint {
	return 0
}
func (bm *BaseModel) SetPK(pk interface{}) error {
	return nil
}
