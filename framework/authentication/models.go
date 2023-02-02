package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"github.com/DmitryVesenniy/go-rest-framework/config"
	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"
	"github.com/DmitryVesenniy/go-rest-framework/framework/typesdb"
)

type User struct {
	serializers.Serializer
	models.BaseModel
	ID              uint                 `json:"id" serializer:"read_only" gorm:"primarykey"`
	Email           string               `json:"email" serializer:"required;title:Емайл;" gorm:"type:varchar(100);unique;unique_index;column:email"`
	Lastname        string               `json:"lastname" serializer:"required;title:Фамилия;" gorm:"column:lastname"`
	Firstname       string               `json:"firstname" serializer:"required;title:Имя;" gorm:"column:firstname"`
	Middlename      string               `json:"middlename" serializer:"required;title:Отчество;" gorm:"column:middlename"`
	Phone           string               `json:"phone" serializer:"title:Телефон;" gorm:"type:varchar(15);column:phone"`
	AdditionalPhone string               `json:"additionalPhone" serializer:"title:Дополнительный номер телефона;" gorm:"type:varchar(15);column:additionalPhone"`
	Avatar          string               `json:"avatar" serializer:"title:Аватар;" gorm:"column:avatar"`
	Password        []byte               `json:"password" serializer:"write_only;title:Пароль;" gorm:"type:bytea;column:password;"`
	IsActive        bool                 `json:"isActive" serializer:"title:Активен;" gorm:"column:isActive;"`
	IsStaff         bool                 `json:"isStaff" serializer:"title:Персонал;" gorm:"column:isStaff;"`
	IsSuperuser     bool                 `json:"isSuperuser" serializer:"title:Суперпользователь;" gorm:"column:isSuperuser;"`
	IsDelete        bool                 `json:"isDelete" serializer:"title:Удален;" gorm:"bit(1);column:isDelete;default:false"`
	RoleID          typesdb.UintNullable `json:"roleId" serializer:"write_only;referenceTable:roles;referencePk:id;" gorm:"column:roleId;"`
	OrganisationID  typesdb.UintNullable `json:"organisationId" serializer:"referenceTable:organisations;referencePk:id;" gorm:"column:organisationId;"`
	CreatedAt       time.Time            `json:"createdAt" serializer:"read_only"`
	UpdatedAt       time.Time            `json:"updatedAt" serializer:"read_only"`
	Role            *Role                `json:"role" serializer:"read_only,title:Роли;referenceField:name;" gorm:"foreignkey:RoleID;references:ID;"`
}

func (User) TableName() string {
	return "users"
}
func (User) Description() string {
	return "Пользователи"
}
func (u *User) Validate() bool {
	isValid, errorData := serializers.GenericValidate(u)
	if len(errorData) > 0 {
		for fileld, errorList := range errorData {
			u.AddError(fileld, errorList...)
		}
	}

	_isValid, _errorList := serializers.ValidateString(
		u.Email,
		serializers.VaidatorRequiredParam(true),
		serializers.ValidatorEmailParam(),
	)
	if !_isValid {
		isValid = false
		u.AddError("email", _errorList...)
	}

	return isValid
}

func (u *User) SetPK(pk interface{}) error {
	id, err := common.ConvertToInt(pk)
	if err != nil {
		return err
	}

	u.ID = uint(id)
	return nil
}
func (u *User) GetPk() uint {
	return u.ID
}

type Role struct {
	serializers.Serializer
	models.BaseModel
	ID                    uint                 `json:"id" serializer:"read_only" gorm:"primarykey"`
	Name                  string               `json:"name" serializer:"required;title:Название роли;" gorm:"column:name;unique;"`
	ModulesPermissionsIDs []uint               `json:"modulesPermissionsIds" serializer:"write_only;referenceTable:module_permissions;referencePk:id;" gorm:"-"`
	ModulesPermissions    []ModulesPermissions `json:"modulesPermissions" serializer:"title:Модули и доступы;referenceField:name&access;" gorm:"many2many:roles_modules;"`
}

func (Role) TableName() string {
	return "roles"
}
func (Role) Description() string {
	return "Роли"
}
func (r *Role) Validate() bool {
	isValid, errorData := serializers.GenericValidate(r)
	if len(errorData) > 0 {
		for fileld, errorList := range errorData {
			r.AddError(fileld, errorList...)
		}
	}
	return isValid
}
func (r *Role) SetPK(pk interface{}) error {
	id, err := common.ConvertToInt(pk)
	if err != nil {
		return err
	}

	r.ID = uint(id)
	return nil
}
func (r *Role) GetPk() uint {
	return r.ID
}

/*
ModulesPermissions - модель описывающая связь модуля и уровнем доступа
Модуль - это лoгический слой, объединяющий в себе несолько таблиц
Модуль и уровни доступа создаются при инициализации приложения,
поэтому в апи не предусмотрены методы CUD для этой таблицы
*/
type ModulesPermissions struct {
	serializers.Serializer
	models.BaseModel
	ID     uint        `json:"id" serializer:"read_only;" gorm:"primarykey"`
	Name   string      `json:"name" serializer:"read_only;title:Название модуля;" gorm:"column:name;uniqueIndex:idx_name_access"`
	Access AccessLevel `json:"access" serializer:"read_only;title:Уровень доступа;" gorm:"column:access;uniqueIndex:idx_name_access"`
	Roles  []Role      `json:"-" gorm:"many2many:roles_modules;"`
}

func (ModulesPermissions) TableName() string {
	return "module_permissions"
}
func (ModulesPermissions) Description() string {
	return "Уровни доступа к модулям"
}

// BaseModelInterface impliments on User
func (u *User) Create(db *gorm.DB) error {
	if u.RoleID != 0 {
		roleDB := &Role{}
		err := db.Where("id = ?", u.RoleID).First(roleDB).Error
		if err != nil {
			return &ErrorNotFoundRole{}
		}
	}

	err := db.Omit("ID", "Role", "Password", "Organisation").Create(u).Error

	return err
}

func (u *User) Update(db *gorm.DB) error {

	if u.RoleID != 0 {
		roleDB := &Role{}
		err := db.Where("id = ?", u.RoleID).First(roleDB).Error
		if err != nil {
			return &ErrorNotFoundRole{}
		}
	}

	err := db.Model(u).
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{
			"email":          u.Email,
			"phone":          u.Phone,
			"firstname":      u.Firstname,
			"lastname":       u.Lastname,
			"middlename":     u.Middlename,
			"isActive":       u.IsActive,
			"isStaff":        u.IsStaff,
			"isSuperuser":    u.IsSuperuser,
			"isDelete":       u.IsDelete,
			"organisationId": u.OrganisationID,
			"roleId":         u.RoleID,
		}).Error

	return err
}

func (u *User) Delete(db *gorm.DB) error {
	u.IsDelete = true
	u.IsActive = false
	return u.Update(db)
}

func (u *User) Clean(db *gorm.DB) error {
	tx := db.Where("id = ?", u.ID).Delete(u)
	return tx.Error
}

func (u *User) SetPassword(password string, db *gorm.DB) error {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("set password: %w", err)
	}

	tx := db.Model(u).
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{
			"password": passwordHashed,
		})

	u.Password = nil
	return tx.Error
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))

	return err == nil
}

func (u *User) Token() (string, error) {
	conf := config.Get()
	var mySigningKey = []byte(conf.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = u.ID

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("something went wrong: %w", err)
	}
	return tokenString, nil
}

func (u *User) RefreshToken() (string, error) {
	conf := config.Get()
	var mySigningKey = []byte(conf.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(conf.RefreshTokenExpaires).Unix()

	refreshTokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("something went wrong: %w", err)
	}
	return refreshTokenString, nil
}

// BaseModelInterface impliments on Role
func (r *Role) Create(db *gorm.DB) error {
	tx := db.Omit("ID", "ModulesPermissions").Create(r)
	if tx.Error != nil {
		return tx.Error
	}

	if len(r.ModulesPermissionsIDs) > 0 {
		modulePermissions := make([]ModulesPermissions, 0)
		err := db.Model(&ModulesPermissions{}).Where("id IN ?", r.ModulesPermissionsIDs).Find(&modulePermissions).Error
		if err != nil {
			return err
		}
		for _, module := range modulePermissions {
			db.Model(r).Association("ModulesPermissions").Append(&module)
		}
	}
	return nil
}

func (r *Role) Update(db *gorm.DB) error {

	err := db.Omit("ID", "ModulesPermissions").Updates(r).Error
	if err != nil {
		return err
	}

	db.Model(r).Association("ModulesPermissions").Clear()

	if len(r.ModulesPermissionsIDs) > 0 {
		modulePermissions := make([]ModulesPermissions, 0)
		err := db.Model(&ModulesPermissions{}).Where("id IN ?", r.ModulesPermissionsIDs).Find(&modulePermissions).Error
		if err != nil {
			return err
		}
		for _, module := range modulePermissions {
			db.Model(r).Association("ModulesPermissions").Append(&module)
		}
	}
	return nil
}

func (r *Role) Delete(db *gorm.DB) error {
	tx := db.Delete(r)
	return tx.Error
}
