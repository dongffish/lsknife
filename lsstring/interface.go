package lsstring

import (
	"fmt"
	"strconv"
)

func InterfaceToString(v interface{}) (string, error) {
	switch val := v.(type) {
	case nil:
		return "", nil
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case int:
		return strconv.Itoa(val), nil
	case int8:
		return strconv.FormatInt(int64(val), 10), nil
	case int16:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case uintptr:
		return strconv.FormatUint(uint64(val), 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(val), 'f', -1, 64), nil // Go 1.15+ 有 FormatComplex
	case complex128:
		return strconv.FormatComplex(val, 'f', -1, 128), nil
	case bool:
		return strconv.FormatBool(val), nil
	case []interface{}:
		// 如果需要自定义格式，可以用 fmt.Sprint
		return fmt.Sprint(val), nil
	default:
		// 对于其他类型，使用默认格式
		fmt.Printf("未知数据类型，使用默认格式： %T\n", v)
		return fmt.Sprint(val), nil
	}
}
