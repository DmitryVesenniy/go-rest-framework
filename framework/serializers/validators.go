package serializers

import (
	"errors"
	"math"
	"reflect"
	"strings"
)

type ErrorMessage string

const (
	StringLengthExceeded ErrorMessage = "String length exceeded"
)

type ValidateFuncOption func(vp *ValidateParams)
type ValidateParams struct {
	Max       int
	Min       int
	MaxLength int
	MinLength int
	Required  bool
	Functions []func(interface{}) (bool, string)
}

func VaidatorRequiredParam(required bool) func(*ValidateParams) {
	return func(vp *ValidateParams) {
		vp.Required = required
	}
}

func VaidatorLengthParam(minLength int, maxLength int) func(*ValidateParams) {
	return func(vp *ValidateParams) {
		vp.MinLength = minLength
		vp.MaxLength = maxLength
	}
}

func ValidatorEmailParam() func(*ValidateParams) {
	return func(vp *ValidateParams) {
		vp.Functions = append(vp.Functions, func(value interface{}) (bool, string) {
			email := value.(string)
			if !strings.Contains(email, "@") {
				return false, "Not valid email"
			}
			return true, ""
		})
	}
}

func ValidateString(val string, opts ...ValidateFuncOption) (bool, []string) {
	errorList := make([]string, 0)
	isValid := true

	vp := &ValidateParams{
		MaxLength: 1000,
		MinLength: 0,
		Required:  false,
	}

	for _, opt := range opts {
		opt(vp)
	}

	if vp.Required && val == "" {
		errorList = append(errorList, "Required parameter")
		isValid = false
	}
	if len(val) < vp.MinLength {
		errorList = append(errorList, "Value is too short")
		isValid = false
	}
	if len(val) > vp.MaxLength {
		errorList = append(errorList, "Value is too long")
		isValid = false
	}

	for _, castomValidFunc := range vp.Functions {
		_isValid, message := castomValidFunc(val)
		if !_isValid {
			isValid = false
			errorList = append(errorList, message)
		}
	}

	return isValid, errorList
}

func ValidateInt(val int, opts ...ValidateFuncOption) (bool, []string) {
	errorList := make([]string, 0)
	isValid := true

	vp := &ValidateParams{
		Min: 0,
		Max: math.MaxInt,
	}

	for _, opt := range opts {
		opt(vp)
	}

	if val < vp.Min {
		errorList = append(errorList, "Too small value")
		isValid = false
	}

	if val > vp.Max {
		errorList = append(errorList, "Too long value")
		isValid = false
	}

	for _, castomValidFunc := range vp.Functions {
		_isValid, message := castomValidFunc(val)
		if !_isValid {
			isValid = false
			errorList = append(errorList, message)
		}
	}

	return isValid, errorList
}

func GenericValidate(serializer SerializersInterface) (bool, map[string][]string) {
	errorList := make(map[string][]string)
	isValid := true

	serializerReflect := reflect.ValueOf(serializer)
	if serializerReflect.Kind() == reflect.Ptr {
		serializerReflect = serializerReflect.Elem()
	}
	for i := 0; i < serializerReflect.NumField(); i++ {
		serializerValueField := serializerReflect.Field(i)
		serializerTypeField := serializerReflect.Type().Field(i)
		serializerFieldName := serializerReflect.Type().Field(i).Name

		tag := serializerTypeField.Tag.Get("serializer")
		tagJSON := serializerTypeField.Tag.Get("json")
		if tag == "" {
			continue
		}

		if tagJSON != "" {
			splited := strings.Split(tagJSON, ",")
			serializerFieldName = splited[0]
		}

		if strings.Contains(tag, SerializerKeyRequired) && serializerValueField.IsZero() {
			errorList[serializerFieldName] = []string{
				"Required parameter",
			}
			isValid = false
		}
		if strings.Contains(tag, SerializerKeyReadOnly) {
			setEmptyValue(serializerValueField)
		}
	}

	return isValid, errorList
}

func setEmptyValue(v reflect.Value) (reflect.Value, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !v.IsValid() {
		return v, errors.New("not valid reflect.Value")
	}

	zero := reflect.Zero(v.Type())
	if !v.CanSet() {
		return zero, errors.New("no CanSet")
	}
	v.Set(zero)
	return zero, nil
}
