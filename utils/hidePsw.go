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

func ToHidePswdString(obj interface{}) (jsonStr string) {
	// 获取传入对象的反射值和类型
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // 解引用
	}
	if val.Kind() != reflect.Struct {
		bytes, _ := json.Marshal(obj)
		return string(bytes)
	}

	typ := val.Type()

	// 创建一个 map 用于存储字段名和值
	output := make(map[string]interface{})

	// 遍历结构体的字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 检查字段是否有 psw 标签
		pswTag := fieldType.Tag.Get("psw")
		if pswTag == "true" && field.Kind() == reflect.String {
			// 处理带有 psw 标签的字段，隐藏密码
			output[fieldType.Name] = maskPassword(field.String())
		} else {
			// 对于非 psw 标签的字段，直接加入输出
			output[fieldType.Name] = field.Interface()
		}
	}

	// 将结果转换为 JSON
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return "error: failed to marshal config to json string"
	}

	return string(jsonBytes)
}
