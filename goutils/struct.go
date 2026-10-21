package goutils

import (
	"reflect"
	"slices"
	"strings"
)

// Struct2Map 将结构体转成map，除去在exclude中声明的tag
func Struct2Map[T any](obj T, exclude ...string) (map[string]any, map[string]any) {
	result := make(map[string]any)
	excluded := make(map[string]any)

	objValue := reflect.Indirect(reflect.ValueOf(obj))
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return result, excluded
	}
	objType := objValue.Type()

	exclude = append(exclude, "-")
	for i := range objValue.NumField() {
		field := objType.Field(i)
		if field.PkgPath != "" {
			// 未导出字段不能读取，跳过
			continue
		}
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "" {
			name = field.Name
		}
		if slices.Contains(exclude, name) {
			excluded[name] = objValue.Field(i).Interface()
		} else {
			result[name] = objValue.Field(i).Interface()
		}
	}
	return result, excluded
}
