package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

func maskPassword(password string) string {
	if len(password) <= 2 {
		return password // 如果密码长度小于等于2，直接返回
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}

func _toHidePswdString(obj interface{}) (obj1 interface{}) {
	val := reflect.ValueOf(obj)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	var valKing = val.Kind()
	//convert slice or array to []interface{}
	if valKing == reflect.Slice || valKing == reflect.Array {
		output := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			output[i] = _toHidePswdString(val.Index(i).Interface())
		}
		return output
	} else if val.Kind() != reflect.Struct {
		return val.Interface()
	}

	typ := val.Type()
	// convert struct to map[string]interface{}
	output := make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		field, fieldType := val.Field(i), typ.Field(i)
		if fieldType.Tag.Get("psw") == "true" && field.Kind() == reflect.String {
			output[fieldType.Name] = maskPassword(field.String())
		} else {
			output[fieldType.Name] = field.Interface()
		}
	}
	return output
}
func ToHidePswdString(obj interface{}) (jsonStr string) {
	o1 := _toHidePswdString(obj)
	// 将结果转换为 JSON
	jsonBytes, err := json.Marshal(o1)
	if err != nil {
		return "error: failed to marshal config to json string"
	}

	return string(jsonBytes)
}
