package difference

import (
	"errors"
	"reflect"
)

type Difference struct {
	OldData interface{}
	NewData interface{}
	ID      interface{}
	data    map[string]interface{}
	err     map[string]interface{}
}

func (d *Difference) CalcDifference() error {
	serializerReflectOld := reflect.ValueOf(d.OldData)
	if serializerReflectOld.Kind() == reflect.Ptr {
		serializerReflectOld = serializerReflectOld.Elem()
	}

	serializerReflecNew := reflect.ValueOf(d.NewData)
	if serializerReflecNew.Kind() == reflect.Ptr {
		serializerReflecNew = serializerReflecNew.Elem()
	}

	if serializerReflectOld.Type().Name() != serializerReflecNew.Type().Name() {
		d.err = map[string]interface{}{
			"error": "type OldData not equals type NewData",
		}
		return errors.New("type OldData not equals type NewData")
	}

	for i := 0; i < serializerReflectOld.NumField(); i++ {
		serializerOldValueField := serializerReflectOld.Field(i)
		serializerOldTypeField := serializerReflectOld.Type().Field(i)
		serializerOldFieldName := serializerReflectOld.Type().Field(i).Name

		serializerNewValueField := serializerReflecNew.Field(i)

		tag := serializerOldTypeField.Tag.Get("serializer")
		if tag == "" {
			continue
		}

		tagJson := serializerOldTypeField.Tag.Get("json")
		if tagJson != "" {
			serializerOldFieldName = tagJson
		}

		kind := serializerOldValueField.Kind()
		if kind == reflect.Struct || kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice {
			continue
		}

		if serializerOldValueField.Interface() != serializerNewValueField.Interface() {
			d.data[serializerOldFieldName] = map[string]interface{}{
				"old": serializerOldValueField.Interface(),
				"new": serializerNewValueField.Interface(),
			}
		}
	}
	return nil
}

func (d *Difference) ToDict() map[string]interface{} {
	return d.data
}
