package lsstring

import (
	"bytes"
	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"unicode/utf8"
)

// DetectTextFileCharset to detext the char set of text files
func DetectTextFileCharset(filename string) (string, error) {
	textDetector := chardet.NewTextDetector()
	buffer := make([]byte, 32<<10)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer f.Close()
	size, _ := io.ReadFull(f, buffer)
	input := buffer[:size]
	var detector = textDetector
	result, err := detector.DetectBest(input)
	if err != nil {
		return "", err
	}
	return result.Charset, nil
}

// DecodeTextToUTF8 test charset of text file and decode non-utf8 to utf8
func DecodeTextToUTF8(filename string) (string, error) {
	textDetector := chardet.NewTextDetector()
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	data := make([]byte, fi.Size())
	size, err := io.ReadFull(f, data)
	if err != nil {
		return "", err
	}
	input := data[:size]
	var detector = textDetector
	result, err := detector.DetectBest(input)
	if err != nil {
		return "", err
	}
	rText := string(input)
	if result.Charset != "UTF-8" {
		rText = mahonia.NewDecoder(result.Charset).ConvertString(rText)
	}
	return rText, nil
}

// ConvertString 字符串编码转换 从srcCharset -> tagCharset。比如 s := ConvertString(s, "gbk", "utf-8")
func ConvertString(src string, srcCharset string, tagCharset string) string {
	if srcCharset == tagCharset {
		return src
	}
	srcCoder := mahonia.NewDecoder(srcCharset)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCharset)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// GuessCharset 猜输入字符串最有可能的编码集
func GuessCharset(src string) (string, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest([]byte(src))
	if err != nil {
		return "", err
	}
	return result.Charset, nil
}

// GBKToUTF8 将数据从BGK转换成UTF8 from bgk to utf8 encoding conversion
func GBKToUTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// GBKToUTF8Str 字符串编码 gbk到utf8转换
func GBKToUTF8Str(str string) (string, error) {
	data, err := GBKToUTF8([]byte(str))
	return string(data), err
}

// UTF8ToGBK utf8 to gbk encoding conversion
func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF8ToGBKStr 字符串编码 utf8到gbk转换
func UTF8ToGBKStr(str string) (string, error) {
	data, err := UTF8ToGBK([]byte(str))
	return string(data), err
}

// IsGBKData 判断数据是否是gbk编码
func IsGBKData(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		// // ASCII 编码的范围:  十进制 => 0 - 127 。  十六进制： 0x00  -  0x7F 。
		if data[i] <= 127 {
			i++
			continue
		} else {
			// GB2312编码的范围: 十进制 => 高位字节：161 - 247, 十六进制：0xA1 - 0xF7
			// 低位字节：161 - 254 , 十六进制：0xA1 - 0xFE
			if data[i] >= 129 &&
				data[i] <= 254 &&
				data[i+1] >= 64 &&
				data[i+1] <= 254 &&
				data[i+1] <= 247 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

// IsGBKStr 判断字符串是否是gbk编码
func IsGBKStr(str string) bool {
	if str == "" {
		return false
	}
	return IsGBKData([]byte(str))
}

// IsUTF8Str 判断字符串是不是UTF8编码
func IsUTF8Str(str string) bool {
	return utf8.ValidString(str)
}
