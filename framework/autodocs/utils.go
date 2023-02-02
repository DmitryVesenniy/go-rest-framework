package autodocs

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/DmitryVesenniy/go-rest-framework/framework/serializers"
)

const (
	MAX_LEVEL = 5
)

var (
	PatternTitle, _          = regexp.Compile(`(title\:(.+?);)|(title\:(.+?)$)`)
	PatternReferenceField, _ = regexp.Compile(`(referenceField\:(.+?);)|(referenceField\:(.+?)$)`)
	PatternReferenceTable, _ = regexp.Compile(`(referenceTable\:(.+?);)|(referenceTable\:(.+?)$)`)
	PatternReferencePk, _    = regexp.Compile(`(referencePk\:(.+?);)|(referencePk\:(.+?)$)`)
	PatternReferenceProxy, _ = regexp.Compile(`(referenceProxy\:(.+?);)|(referenceProxy\:(.+?)$)`)
)

func GetSerializerFields(serializer interface{}, level int) []SerializerField {
	if level > MAX_LEVEL {
		return nil
	}
	serializerReflect := reflect.ValueOf(serializer)
	if serializerReflect.Kind() == reflect.Ptr {
		serializerReflect = serializerReflect.Elem()
	}

	serializersField := make([]SerializerField, 0)
	switch serializerReflect.Kind() {
	case reflect.Array, reflect.Slice:
		arrReflect := serializerReflect.Type().Elem()
		reflectionValue := reflect.New(arrReflect)
		serializerField := SerializerField{
			NameField: serializerReflect.Type().Name(),
			Type:      "array",
		}
		serializerField.Signatura = GetSerializerFields(reflectionValue.Interface(), level+1)
		serializersField = append(serializersField, serializerField)
	case reflect.Struct:
		for i := 0; i < serializerReflect.NumField(); i++ {
			typeField := serializerReflect.Type().Field(i)
			field := serializerReflect.Field(i)
			fieldKind := field.Kind()

			typeName := ""

			if typeField.Type.Kind() == reflect.Ptr {
				typeName = typeField.Type.Elem().Name()
			} else {
				typeName = typeField.Type.Name()
			}

			if field.Kind() == reflect.Ptr {
				fieldKind = field.Elem().Kind()
			}

			tag := typeField.Tag.Get("serializer")
			tagJSON := typeField.Tag.Get("json")

			if tag == "" {
				continue
			}

			serializerFieldName := typeField.Name
			if tagJSON != "-" {
				splited := strings.Split(tagJSON, ",")
				if len(splited) > 0 {
					serializerFieldName = splited[0]
				}
			}

			var signatura interface{} = nil
			if !field.IsValid() {
				continue
			}
			_typeField := field.Type().Name()

			var referenceData *ReferenceData = nil
			var referenceProxy *ReferenceData = nil

			proxyField := PatternReferenceProxy.FindString(tag)
			if proxyField != "" {
				referenceProxy = &ReferenceData{
					PresentationField: strings.Trim(proxyField[15:], ";"),
				}
			}

			presetationField := PatternReferenceField.FindString(tag)
			if presetationField != "" {
				referenceData = &ReferenceData{
					PresentationField: strings.Trim(presetationField[15:], ";"),
				}
			}

			switch fieldKind {
			case reflect.Slice, reflect.Array:
				_typeField = "array"
				arrReflect := field.Type().Elem()
				if arrReflect.Name() == serializerReflect.Type().Name() {
					// убираем рекурсию, если елемент вложенный в массив имеет такой же тип, что и родительский элемент
					signatura = "self"
				} else {
					reflectionValue := reflect.New(arrReflect)
					signatura = GetSerializerFields(reflectionValue.Interface(), level+1)
				}
			case reflect.Struct:
				if typeName == "Time" {
					_typeField = "date"
				} else {
					_typeField = "object"
					referenseObj := field.Interface()
					signatura = GetSerializerFields(referenseObj, level+1)
				}
			}

			if typeName == "Time" {
				_typeField = "date"
			}

			var referenseTabe *ReferenceTable = nil
			tableName := PatternReferenceTable.FindString(tag)
			pkField := PatternReferencePk.FindString(tag)
			if tableName != "" || pkField != "" {
				referenseTabe = &ReferenceTable{
					TableName: strings.Trim(tableName[15:], ";"),
					PkField:   strings.Trim(pkField[12:], ";"),
				}
			}

			serializerField := SerializerField{
				NameField:      serializerFieldName,
				Type:           _typeField,
				Signatura:      signatura,
				ReferenceData:  referenceData,
				ReferenceTable: referenseTabe,
				ProxyField:     referenceProxy,
			}
			if strings.Contains(tag, serializers.SerializerKeyRequired) {
				serializerField.Required = true
			}
			if strings.Contains(tag, serializers.SerializerKeyWriteOnly) {
				serializerField.WriteOnly = true
			}
			if strings.Contains(tag, serializers.SerializerKeyReadOnly) {
				serializerField.ReadOnly = true
			}
			title := PatternTitle.FindString(tag)
			if title != "" {
				serializerField.Title = strings.Trim(title[6:], ";")
			}
			serializersField = append(serializersField, serializerField)
		}
	}
	return serializersField
}
