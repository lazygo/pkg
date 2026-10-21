package goutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// String 转换为字符串
// 如果value不能转换为字符串，则返回默认值defVal
// defVal 提供默认值，如果没有，则视为空字符串""
func String(value any, defVal ...string) string {
	var str string
	switch value := value.(type) {
	case string:
		str = value
	case []byte:
		str = string(value)
	case int, int8, int16, int32, int64:
		str = strconv.FormatInt(reflect.ValueOf(value).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		str = strconv.FormatUint(reflect.ValueOf(value).Uint(), 10)
	case float32:
		str = strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64)
	default:
		str = ""
	}

	if str == "" && len(defVal) == 1 {
		return defVal[0]
	}
	return str
}

// Int64 转换为int64类型
// value 可以是数字、数字字符串
// 如果value不能转换为数字，则返回默认值defVal
// defVal 提供默认值，如果没有，则视为0
func Int64(value any, defVal ...int64) int64 {
	val := reflect.ValueOf(value)
	var d int64
	var err error
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	case float32, float64:
		d = int64(val.Float())
	case string:
		d, err = strconv.ParseInt(strings.TrimSpace(val.String()), 10, 64)
	case []byte:
		d, err = strconv.ParseInt(strings.TrimSpace(string(value)), 10, 64)
	default:
		err = fmt.Errorf("Int64 need numeric not `%T`", value)
	}
	if err != nil {
		if len(defVal) == 1 {
			return defVal[0]
		}
		return 0
	}
	return d
}

// Uint64 转换为uint64类型
// value 可以是数字、数字字符串
// 如果value不能转换为数字，则返回默认值defVal
// defVal 提供默认值，如果没有，则视为0
func Uint64(value any, defVal ...uint64) uint64 {
	val := reflect.ValueOf(value)
	var d uint64
	var err error
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		d = uint64(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = val.Uint()
	case float32, float64:
		d = uint64(val.Float())
	case string:
		d, err = strconv.ParseUint(strings.TrimSpace(val.String()), 10, 64)
	case []byte:
		d, err = strconv.ParseUint(strings.TrimSpace(string(value)), 10, 64)
	default:
		err = fmt.Errorf("Uint64 need numeric not `%T`", value)
	}
	if err != nil {
		if len(defVal) == 1 {
			return defVal[0]
		}
		return 0
	}
	return d
}

// Int 转换为int类型
// value 可以是数字、数字字符串
// 如果value不能转换为数字，则返回默认值defVal
// defVal 提供默认值，如果没有，则视为0
func Int(value any, defVal ...int) int {
	val := reflect.ValueOf(value)
	var d int
	var err error
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		d = int(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = int(val.Uint())
	case float32, float64:
		d = int(val.Float())
	case string:
		d, err = strconv.Atoi(strings.TrimSpace(val.String()))
	case []byte:
		d, err = strconv.Atoi(strings.TrimSpace(string(value)))
	default:
		err = fmt.Errorf("Int need numeric not `%T`", value)
	}
	if err != nil {
		if len(defVal) == 1 {
			return defVal[0]
		}
		return 0
	}

	return d
}

// Uint 转换为uint类型
// value 可以是数字、数字字符串
// 如果value不能转换为数字，则返回默认值defVal
// defVal 提供默认值，如果没有，则视为0
func Uint(value any, defVal ...uint) uint {
	val := reflect.ValueOf(value)
	var d uint
	var err error
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		d = uint(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = uint(val.Uint())
	case float32, float64:
		d = uint(val.Float())
	case string:
		var d64 uint64
		d64, err = strconv.ParseUint(strings.TrimSpace(val.String()), 10, strconv.IntSize)
		d = uint(d64)
	case []byte:
		var d64 uint64
		d64, err = strconv.ParseUint(strings.TrimSpace(string(value)), 10, strconv.IntSize)
		d = uint(d64)
	default:
		err = fmt.Errorf("Uint need numeric not `%T`", value)
	}
	if err != nil {
		if len(defVal) == 1 {
			return defVal[0]
		}
		return 0
	}

	return d
}

// Float 转换为float64类型
// value 可以是数字、数字字符串
// 如果value不能转换为数字，则返回默认值defVal
// defVal 提供默认值，如果没有，则视为0
func Float(value any, defVal ...float64) float64 {
	val := reflect.ValueOf(value)
	var d float64
	var err error
	switch value := value.(type) {
	case int, int8, int16, int32, int64:
		d = float64(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		d = float64(val.Uint())
	case float32, float64:
		d = val.Float()
	case string:
		d, err = strconv.ParseFloat(strings.TrimSpace(val.String()), 64)
	case []byte:
		d, err = strconv.ParseFloat(strings.TrimSpace(string(value)), 64)
	default:
		err = fmt.Errorf("Float need numeric not `%T`", value)
	}
	if err != nil {
		if len(defVal) == 1 {
			return defVal[0]
		}
		return 0
	}

	return d
}

// Bool 真值测试
func Bool(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case string:
		return v != ""
	case []any:
		return len(v) != 0
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return !val.IsZero()
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() > 0
	default:
		return val.IsValid() && !val.IsZero()
	}
}
