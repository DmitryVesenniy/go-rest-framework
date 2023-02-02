package common

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var (
	PhoneRegexp = regexp.MustCompile(`[a-zA-Z\s-\(\)\,\:]`)
)

// ConvertToInt конвертируем интерфейс в int
func ConvertToInt(data interface{}) (int, error) {
	switch v := data.(type) {
	case float64:
		value := data.(float64)
		return int(value), nil
	case float32:
		value := data.(float32)
		return int(value), nil
	case int:
		value := data.(int)
		return int(value), nil
	case int64:
		value := data.(int64)
		return int(value), nil
	case uint:
		value := data.(uint)
		return int(value), nil
	case uint64:
		value := data.(uint64)
		return int(value), nil
	case uint32:
		value := data.(uint32)
		return int(value), nil
	case json.Number:
		res := data.(json.Number).String()
		value, err := strconv.Atoi(res)
		return int(value), err
	case string:
		res := data.(string)
		value, err := strconv.Atoi(res)
		return int(value), err
	case []byte:
		value := binary.BigEndian.Uint64(v)
		return int(value), nil
	}
	return 0, errors.New("error convert intrface{} to int")
}

// ConvertToString конвертируем interface{} в строку
func ConvertToString(data interface{}) (res string) {
	switch v := data.(type) {
	case float64:
		res = strconv.FormatFloat(data.(float64), 'f', 6, 64)
	case float32:
		res = strconv.FormatFloat(float64(data.(float32)), 'f', 6, 32)
	case int:
		res = strconv.FormatInt(int64(data.(int)), 10)
	case int64:
		res = strconv.FormatInt(data.(int64), 10)
	case uint:
		res = strconv.FormatUint(uint64(data.(uint)), 10)
	case uint64:
		res = strconv.FormatUint(data.(uint64), 10)
	case uint32:
		res = strconv.FormatUint(uint64(data.(uint32)), 10)
	case json.Number:
		res = data.(json.Number).String()
	case string:
		res = data.(string)
	case []byte:
		res = string(v)
	case time.Time:
		res = fmt.Sprintf("%d-%d-%d", v.Day(), v.Month(), v.Year())
	default:
		res = ""
	}
	return
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func PhoneConverter(phone string) string {
	return PhoneRegexp.ReplaceAllString(phone, "")
}

func ConvertBoolToStr(v bool, associated map[bool]string) string {
	if associated == nil {
		associated = map[bool]string{
			true:  "Да",
			false: "Нет",
		}
	}

	return associated[v]
}
