package serializers

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/DmitryVesenniy/go-rest-framework/framework/models"
)

type Serializer struct {
	fields        map[string]Rules
	errors        map[string][]string
	InstanceModel models.BaseModelInterface `json:"-" gorm:"-"`
}

func (s *Serializer) Validate() bool {
	return true
}

func (s *Serializer) Data() map[string]interface{} {
	return map[string]interface{}{}
}

func (s *Serializer) SetError(err string) {
	s.errors = map[string][]string{"details": {err}}
}

func (s *Serializer) AddError(field string, err ...string) {
	if s.errors == nil {
		s.errors = make(map[string][]string)
	}
	s.errors[field] = append(s.errors[field], err...)
}

func (s *Serializer) Errors() map[string][]string {
	return s.errors
}

func (s *Serializer) GetModel() models.BaseModelInterface {
	return s.InstanceModel
}

func (s *Serializer) GetInstanceModel() models.BaseModelInterface {
	return s.InstanceModel
}

func InitValuesSerializer(serializer SerializersInterface, model models.BaseModelInterface) SerializersInterface {
	serializerReflect := reflect.ValueOf(serializer)
	if serializerReflect.Kind() == reflect.Ptr {
		serializerReflect = serializerReflect.Elem()
	}
	modelReflect := reflect.ValueOf(model)
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	for i := 0; i < serializerReflect.NumField(); i++ {
		serializerValueField := serializerReflect.Field(i)
		serializerTypeField := serializerReflect.Type().Field(i)
		serializerFieldName := serializerReflect.Type().Field(i).Name

		// if serializerValueField.Kind() == reflect.Ptr {
		// 	serializerValueField = serializerValueField.Elem()
		// }

		tag := serializerTypeField.Tag.Get("serializer")
		if tag == "" || strings.Contains(tag, SerializerKeyWriteOnly) {
			continue
		}

		for i := 0; i < modelReflect.NumField(); i++ {
			modelValueField := modelReflect.Field(i)
			modelTypeField := modelReflect.Type().Field(i)
			modelFieldName := modelReflect.Type().Field(i).Name

			// if modelValueField.Kind() == reflect.Ptr {
			// 	modelValueField = modelValueField.Elem()
			// }

			if serializerFieldName == modelFieldName && serializerTypeField.Type.Name() == modelTypeField.Type.Name() {
				serializerValueField.Set(modelValueField)
			}
		}
	}

	return serializer
}

func InitValuesModel(serializer SerializersInterface) models.BaseModelInterface {
	model := serializer.GetModel()
	if model == nil {
		return nil
	}

	serializerReflect := reflect.ValueOf(serializer)
	modelReflect := reflect.ValueOf(model)
	if serializerReflect.Kind() == reflect.Ptr {
		serializerReflect = serializerReflect.Elem()
	}
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	for i := 0; i < serializerReflect.NumField(); i++ {
		valueField := serializerReflect.Field(i)
		typeField := serializerReflect.Type().Field(i)
		fieldName := serializerReflect.Type().Field(i).Name

		tag := typeField.Tag.Get("serializer")
		if tag == "" || strings.Contains(tag, SerializerKeyReadOnly) {
			continue
		}

		for i := 0; i < modelReflect.NumField(); i++ {
			modelReflectValueField := modelReflect.Field(i)
			modelReflectTypeField := modelReflect.Type().Field(i)
			modelReflectFieldName := modelReflect.Type().Field(i).Name

			if fieldName == modelReflectFieldName && typeField.Type.Name() == modelReflectTypeField.Type.Name() {
				modelReflectValueField.Set(valueField)
			}
		}
	}
	return model
}

func MarshallSerializer(serializer interface{}) ([]byte, error) {
	serializer = fitSerializerData(serializer)
	b, err := json.Marshal(serializer)
	return b, err
}

func FitSerializerData(serializer interface{}) interface{} {
	return fitSerializerData(serializer)
}

func fitSerializerData(serializer interface{}) interface{} {
	serializerReflect := reflect.ValueOf(serializer)
	if serializerReflect.Kind() == reflect.Ptr {
		serializerReflect = serializerReflect.Elem()
	}

	switch serializerReflect.Kind() {
	case reflect.Map:
		// mv, ok := serializer.(map[string]interface{})
		// if ok {
		// 	for key := range mv {
		// 		mapValueRef := reflect.ValueOf(mv[key])
		// 		if mapValueRef.Kind() == reflect.Ptr {
		// 			mapValueRef = mapValueRef.Elem()
		// 		}
		// 		mv[key] = fitSerializerData(mapValueRef.Interface())
		// 	}
		// }
		// return mv
		for _, k := range serializerReflect.MapKeys() {
			item := serializerReflect.MapIndex(k)
			fitRefData(item.Elem())
		}
	case reflect.Array, reflect.Slice:
		if serializerReflect.Len() > 0 {
			for i := 0; i < serializerReflect.Len(); i++ {
				elem := serializerReflect.Index(i)
				fitRefData(elem)
			}
		}
	case reflect.Struct:
		for i := 0; i < serializerReflect.NumField(); i++ {
			serializerValueField := serializerReflect.Field(i)
			serializerTypeField := serializerReflect.Type().Field(i)

			tag := serializerTypeField.Tag.Get("serializer")

			if strings.Contains(tag, SerializerKeyWriteOnly) {
				setEmptyValue(serializerValueField)
			}
		}
	}
	return serializer
}

func fitRefData(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			item := v.MapIndex(k)
			fitRefData(item)
		}
	case reflect.Array, reflect.Slice:
		if v.Len() > 0 {
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				fitRefData(elem)
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			serializerValueField := v.Field(i)
			serializerTypeField := v.Type().Field(i)

			tag := serializerTypeField.Tag.Get("serializer")

			if strings.Contains(tag, SerializerKeyWriteOnly) {
				setEmptyValue(serializerValueField)
			} else {
				if serializerValueField.Kind() == reflect.Struct || serializerValueField.Kind() == reflect.Slice || serializerValueField.Kind() == reflect.Map {
					fitRefData(serializerValueField)
				}
			}
		}
	}
}
