package authentication

import (
	"fmt"
	"reflect"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"gorm.io/gorm"
)

var (
	UserConvertersData = map[string]func(*User, interface{}, *gorm.DB) error{
		"id": func(u *User, i interface{}, db *gorm.DB) error {
			id, err := common.ConvertToInt(i)
			if err != nil {
				return fmt.Errorf("error convert interface{}[%v] to ID(uint): %w", i, err)
			}
			u.ID = uint(id)
			return nil
		},
		"email": func(u *User, i interface{}, db *gorm.DB) error {
			email, ok := i.(string)
			if !ok {
				return fmt.Errorf("convert interface{}[%v] to Email(string)", i)
			}
			u.Email = email
			return nil
		},
		"is_active": func(u *User, i interface{}, db *gorm.DB) error {
			isActive, ok := i.(bool)
			if !ok {
				return fmt.Errorf("convert interface{}[%v] to IsActive(bool)", i)
			}
			u.IsActive = isActive
			return nil
		},
		"is_superuser": func(u *User, i interface{}, db *gorm.DB) error {
			isSuperuser, ok := i.(bool)
			if !ok {
				return fmt.Errorf("convert interface{}[%v] to IsSuperuser(bool)", i)
			}
			u.IsSuperuser = isSuperuser
			return nil
		},
		"role": func(u *User, i interface{}, db *gorm.DB) error {
			if i == nil {
				return nil
			}

			s := reflect.ValueOf(i)
			if s.Kind() != reflect.Map {
				return fmt.Errorf("interfaceStruct() given a non-struct type")
			}
			if s.IsNil() {
				return nil
			}

			roleMap, isOk := i.(map[string]interface{})
			if !isOk {
				return fmt.Errorf("error convert map to role")
			}

			role := &Role{}

			if err := db.Model(&Role{}).Preload("ModulesPermissions").Where("id = ?", roleMap["id"]).First(role).Error; err == nil {
				u.Role = role
			} else {
				return err
			}
			return nil
		},
	}
)

func GetUserFromContext(userData interface{}, db *gorm.DB) (*User, error) {
	if userData == nil {
		return nil, fmt.Errorf("GetUserFromContext: user nil pointer")
	}
	user := &User{}
	userMap := make(map[string]interface{})
	v := reflect.ValueOf(userData)
	if v.IsZero() {
		return nil, fmt.Errorf("GetUserFromContext: user nil pointer")
	}
	if v.Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			key, ok := iter.Key().Interface().(string)
			if !ok {
				continue
			}
			value := iter.Value().Interface()
			userMap[key] = value
		}
	} else {
		return user, fmt.Errorf("interface does not match map format")
	}

	id := int(userMap["id"].(float64))

	if err := db.Model(&User{}).Preload("Role").Preload("Role.ModulesPermissions").Where("id = ?", id).First(user).Error; err != nil {
		return user, err
	}

	return user, nil
}
