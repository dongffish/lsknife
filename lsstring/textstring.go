package lsstring

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// GetStringBetween 获取两个子串之间的内容这个函数GetStringBetween接受三个参数：一个主字符串str，一个起始子串start和一个结束子串end。函数首先查找起始子串和结束子串在主字符串中的位置，然后返回它们之间的内容。如果任一子串未找到，函数返回空字符串。
func GetStringBetween(str, start, end string) string {
	// 使用 strings.Index 查找开始子串的位置
	startIndex := strings.Index(str, start)
	if startIndex < 0 {
		return ""
	}
	// 使用 strings.Index 查找结束子串的位置
	endIndex := strings.Index(str, end)
	if endIndex < 0 || endIndex <= startIndex {
		return ""
	}
	// 使用字符串切片获取两个子串之间的内容
	return str[startIndex+len(start) : endIndex]
}

// RemoveInvisibleChars 删除不可见字符
func RemoveInvisibleChars(input string) string {
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	return re.ReplaceAllString(input, "")
}

// FindStringWithRegexp 根据正则表达是找出所有匹配串
func FindStringWithRegexp(src string, regstr string) ([]string, bool) {
	compiled := regexp.MustCompile(regstr)
	if compiled != nil {
		params := compiled.FindAllString(src, -1)
		return params, len(params) != 0
	}
	return nil, false
}

// FindIPStrings 根据正则表达是找出所有匹配IP串
func FindIPStrings(src string) ([]string, bool) {
	return FindStringWithRegexp(src, ipaddrPattern)
}

// DeleteInvisibleCh 删除不可见字符
func DeleteInvisibleCh(str string) string {
	srcRunes := []rune(str)
	dstRuns := make([]rune, 0, len(srcRunes))
	for _, c := range srcRunes {
		if c >= 0 && c <= 31 {
			continue
		}
		if c == 127 {
			continue
		}
		dstRuns = append(dstRuns, c)
	}
	return string(dstRuns)
}

// FirstToUpper 字符串首字母大写
func FirstToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstToLower 字符串首字母小写
func FirstToLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func InterfaceToString(v interface{}) (string, error) {
	// 使用 type switch 判断 b 的实际类型
	switch v := v.(type) {
	case nil:
		return "", errors.New("返回值为 nil")
	case string:
		return fmt.Sprintf("%s", v), nil
	case []byte:
		return string(v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case float64:
		fmt.Printf("类型是 float64: %.2f\n", v)
	case bool:
		return fmt.Sprintf("%t", v), nil
	case []interface{}:
		//fmt.Printf("类型是 []interface{}: %v\n", v)
		return fmt.Sprintf("%v", v), nil
	}
	return fmt.Sprintf("%v", v), errors.New(fmt.Sprintf("未知数据类型 %T", v))
}
