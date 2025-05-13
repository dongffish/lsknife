package lsstring

import (
	"regexp"
	"unicode"
)

const (
	// 汉字匹配正则表达式
	hanStringPattern = "^[\u4e00-\u9fa5]{3,8}$"
	// 中国大陆手机号码正则匹配, 不是那么太精细
	// 只要是 13,14,15,17,18 开头的 11 位数字就认为是中国手机号
	chinaMobilePattern = `^1[34578][0-9]{9}$`
	// 用户昵称的正则匹配, 合法的字符有 0-9, A-Z, a-z, _, 汉字
	// 字符 '_' 只能出现在中间且不能重复, 如 "__"
	nicknamePattern = `^[a-z0-9A-Z\p{Han}]+(_[a-z0-9A-Z\p{Han}]+)*$`
	// 用户名的正则匹配, 合法的字符有 0-9, A-Z, a-z, _
	// 第一个字母不能为 _, 0-9
	// 最后一个字母不能为 _, 且 _ 不能连续
	usernamePattern = `^[a-zA-Z][a-z0-9A-Z]*(_[a-z0-9A-Z]+)*$`
	// 电子邮箱的正则匹配, 考虑到各个网站的 mail 要求不一样, 这里匹配比较宽松
	// 邮箱用户名可以包含 0-9, A-Z, a-z, -, _, .
	// 开头字母不能是 -, _, .
	// 结尾字母不能是 -, _, .
	// -, _, . 这三个连接字母任意两个不能连续, 如不能出现 --, __, .., -_, -., _.
	// 邮箱的域名可以包含 0-9, A-Z, a-z, -
	// 连接字符 - 只能出现在中间, 不能连续, 如不能 --
	// 支持多级域名, x@y.z, x@y.z.w, x@x.y.z.w.e
	mailPattern = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*\.)+[a-zA-Z]{2,4}$`
	// ip地址校验正则表达式
	//ipaddrPattern = `/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/`
	ipaddrPattern = `((?:(?:25[0-5]|2[0-4]\d|(?:1\d{2}|[1-9]?\d))\.){3}(?:25[0-5]|2[0-4]\d|(?:1\d{2}|[1-9]?\d)))`
	// 验证数字和26个字母组成的字符串
	alphanumberPattern = `^[A-Za-z0-9]+$`
	// 验证数字和26个字母,以及下划线组成的字符串
	alphanumberunderlinePattern = `^[A-Za-z0-9_]+$`
	// 验证数字：^[0-9]*$
	numberPattern = `^[0-9]*$`
	// 验证由26个英文字母组成的字符串：^[A-Za-z]+$
	alphaPattern = `^[A-Za-z]+$`
	// 验证由26个大写英文字母组成的字符串：^[A-Z]+$
	uppercasePattern = `^[A-Z]+$`
	// 验证由26个小写英文字母组成的字符串：^[a-z]+$
	lowercasePattern = `^[a-z]+$`
)

var (
	hanStringRegexp            = regexp.MustCompile(hanStringPattern)
	chinaMobileRegexp          = regexp.MustCompile(chinaMobilePattern)
	nicknameRegexp             = regexp.MustCompile(nicknamePattern)
	usernameRegexp             = regexp.MustCompile(usernamePattern)
	mailRegexp                 = regexp.MustCompile(mailPattern)
	ipAddrRegexp               = regexp.MustCompile(ipaddrPattern)
	alphaRegexp                = regexp.MustCompile(alphaPattern)
	alphanumberRegexp          = regexp.MustCompile(alphanumberPattern)
	numberRegexp               = regexp.MustCompile(numberPattern)
	uppercaseRegexp            = regexp.MustCompile(uppercasePattern)
	lowercaseRegexp            = regexp.MustCompile(lowercasePattern)
	alphanumberunderlineRegexp = regexp.MustCompile(alphanumberunderlinePattern)
)

/*
//判断是不是只有字母
const alpha = "abcdefghijklmnopqrstuvwxyz"
func IsAlphaOnly(s string) bool {
	for _, char := range s {
		if !strings.Contains(alpha, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}

//判断是不是只有数字
const number = "0123456789"

func IsNumberOnly(s string) bool {
	for _, char := range s {
		if !strings.Contains(number, string(char)) {
			return false
		}
	}
	return true
}

//判断是不是只有字母和数字
const alpha_number = "abcdefghijklmnopqrstuvwxyz0123456789"

func IsAlphaNumberOnly(s string) bool {
	for _, char := range s {
		if !strings.Contains(alpha_number, strings.ToLower(string(char))) {
			return false
		}
	}
	return true
}
*/

// IsAlphaNumberOnly 判断数字和26字母组成的串
func IsAlphaNumberOnly(str string) bool {
	return alphanumberRegexp.MatchString(str)
}

// IsAlphaNumberUnderlineOnly 判断数字、下划线和26字母组成的串
func IsAlphaNumberUnderlineOnly(str string) bool {
	return alphanumberunderlineRegexp.MatchString(str)
}

// IsAlphaOnly 判断26个字母组成的串
func IsAlphaOnly(str string) bool {
	return alphaRegexp.MatchString(str)
}

// IsNumberOnly 判断只有数字组成的串
func IsNumberOnly(str string) bool {
	return numberRegexp.MatchString(str)
}

// IsUpperCaseOnly 判断只有大写字母组成的穿
func IsUpperCaseOnly(str string) bool {
	return uppercaseRegexp.MatchString(str)
}

// IsLowerCaseOnly 判断只有小写字母组成的串
func IsLowerCaseOnly(str string) bool {
	return lowercaseRegexp.MatchString(str)
}

// IsChineseString 判断全中文字符串
func IsChineseString(str string) bool {
	return hanStringRegexp.MatchString(str)
}

// IsChineseChar 判断中文串
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// IsChinaMobile 检验是否为合法的中国手机号, 不是那么太精细
// 只要是 13,14,15,18 开头的 11 位数字就认为是中国手机号
func IsChinaMobile(b []byte) bool {
	if len(b) != 11 {
		return false
	}
	return chinaMobileRegexp.Match(b)
}

// IsChinaMobileString 同 func IsChinaMobile(b []byte) bool
func IsChinaMobileString(str string) bool {
	if len(str) != 11 {
		return false
	}
	return chinaMobileRegexp.MatchString(str)
}

// IsNickname 检验是否为合法的昵称, 合法的字符有 0-9, A-Z, a-z, _, 汉字
// 字符 '_' 只能出现在中间且不能重复, 如 "__"
func IsNickname(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return nicknameRegexp.Match(b)
}

// IsNicknameString 同 func IsNickname(b []byte) bool
func IsNicknameString(str string) bool {
	if len(str) == 0 {
		return false
	}
	return nicknameRegexp.MatchString(str)
}

// IsUserName 检验是否为合法的用户名, 合法的字符有 0-9, A-Z, a-z, _
// 第一个字母不能为 _, 0-9
// 最后一个字母不能为 _, 且 _ 不能连续
func IsUserName(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return usernameRegexp.Match(b)
}

// IsUserNameString 同 func IsName(b []byte) bool
func IsUserNameString(str string) bool {
	if len(str) == 0 {
		return false
	}
	return usernameRegexp.MatchString(str)
}

// IsMail 检验是否为合法的电子邮箱, 考虑到各个网站的 mail 要求不一样, 这里匹配比较宽松
// 邮箱用户名可以包含 0-9, A-Z, a-z, -, _, .
// 开头字母不能是 -, _, .
// 结尾字母不能是 -, _, .
// -, _, . 这三个连接字母任意两个不能连续, 如不能出现 --, __, .., -_, -., _.
// 邮箱的域名可以包含 0-9, A-Z, a-z, -
// 连接字符 - 只能出现在中间, 不能连续, 如不能 --
// 支持多级域名, x@y.z, x@y.z.w, x@x.y.z.w.e
func IsMail(b []byte) bool {
	if len(b) < 6 { // x@x.xx
		return false
	}
	return mailRegexp.Match(b)
}

// IsMailString 同 func IsMail(b []byte) bool
func IsMailString(str string) bool {
	if len(str) < 6 { // x@x.xx
		return false
	}
	return mailRegexp.MatchString(str)
}

// IsIPAddress validate ipv4 address
func IsIPAddress(b []byte) bool {
	if len(b) > 15 { // xxx.xxx.xxx.xxx
		return false
	}
	return ipAddrRegexp.Match(b)
}

// IsIPAddressString 同 func IsIPAddress(b []byte) bool
func IsIPAddressString(str string) bool {
	if len(str) > 15 { // xxx.xxx.xxx.xxx
		return false
	}
	return ipAddrRegexp.MatchString(str)
}
