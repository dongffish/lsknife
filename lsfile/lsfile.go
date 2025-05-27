package lsfile

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GetCurrentExecutePath 获取当时运行程序的目录 get current exe path
func GetCurrentExecutePath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetExecuteFileName 获取当时运行程序的名称 get exe file name
func GetExecuteFileName() (string, error) {
	dir := strings.Split(os.Args[0], string(os.PathSeparator))
	filename := dir[len(dir)-1]
	suffix := path.Ext(filename)
	return strings.TrimSuffix(filename, suffix), nil
}

// CopyFile copy files
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// WriteFile 写文件(一次性写入,原文件存时则覆盖内容)
func WriteFile(fileName string, Content string) error {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
	//使用完毕,需要关闭文件
	defer f.Close()

	_, err = f.WriteString(Content)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// WriteFileAppend 写文件(追加写入内容到末尾)
func WriteFileAppend(fileName string, Content string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 打开文件，如果文件不存在则创建
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	/*
		//打开文件,新建文件
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println("err = ", err)
			Logger.Errorln(GetLogPrefix(), err.Error())
			return
		}*/
	//使用完毕,需要关闭文件
	defer f.Close()

	_, err = f.WriteString(Content)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// PathExists test path or file is exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetFileList browse directory
func GetFileList(root string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	return fileList, err
}

// GetFileListWithFilter browse directory
func GetFileListWithFilter(root string, filter []string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(root, func(dir string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		for _, v := range filter {
			if strings.HasSuffix(dir, v) {
				fileList = append(fileList, dir)
				break
			}
		}
		return nil
	})
	return fileList, err
}

// FileSize 获取文件大写
func FileSize(file string) (int64, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// FileInfo 获取文件信息
func FileInfo(file string) (os.FileInfo, error) {
	return os.Stat(file)
}
