package main

import (
	"fmt"
	"gitee.com/dongffish/lsknife/lsencrypt"
	"github.com/voxelbrain/goptions"
)

func main() {
	options := struct {
		Func int           `goptions:"-f, --func, obligatory, description='功能 1-加密 2-解密'"`
		Key  string        `goptions:"-k, --key, obligatory, description='密钥'"`
		Msg  string        `goptions:"-m, --msg, obligatory, description='字符串'"`
		Help goptions.Help `goptions:"-h, --help, description='打印帮助'"`
	}{}

	err := goptions.Parse(&options)
	if err != nil {
		printHelp()
		return
	}

	if options.Func != 1 && options.Func != 2 {
		fmt.Println("请指定功能，1-加密 2-解密")
		return
	}
	if options.Func == 1 { // 加密
		res, err := lsencrypt.TripleDesEncrypt([]byte(options.Msg), []byte(options.Key))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("加密后的字符串为:", string(res))
	}
	if options.Func == 2 {
		res, err := lsencrypt.TripleDesDecrypt([]byte(options.Msg), []byte(options.Key))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("解密后的字符串为:", string(res))
	}
}

func printHelp() {
	fmt.Println("t3des 3des加解密")
	fmt.Println("-----------------------")
	fmt.Println("Example:")
	fmt.Println("  t3des -f 1 -k 1234567 -m abcdefg")
	fmt.Println("  加密后的字符串为: g5Sqlo0EDZE=")
	fmt.Println("Options:")
	fmt.Println("  -f, --func 功能 1-加密 2-解密(*)")
	fmt.Println("  -k, --key 密钥 (*)")
	fmt.Println("  -m, --msg 当功能是加密时，这里放待加密的明文。当功能是解密时，这里放待解密的密文")
	fmt.Println("  -h, --help      打印帮助")
}
