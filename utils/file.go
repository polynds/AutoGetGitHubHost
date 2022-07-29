package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func SaveFile(filePath string, data io.ReadCloser) {
	if CheckAndCreateDir(filePath) {
		fmt.Println("文件夹创建失败")
		return
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("保存文件打开失败", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("保存文件关闭失败", err)
		}
	}(file)

	_, err = io.Copy(file, data)
	if err != nil {
		fmt.Println("保存文件写入失败", err)
	}
}

func CheckAndCreateDir(filePath string) bool {
	dir := path.Dir(filePath)
	if !IsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("保存文件创建失败", err)
			return true
		}
	}
	return false
}

func ReadAll(filePath string) []byte {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("读取文件打开失败", err)
		return nil
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("读取文件关闭失败", err)
		}
	}(file)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	return data
}

func WriteAll(filePath string, data string) error {

	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("写入文件打开失败", err)
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("写入文件关闭失败", err)
		}
	}(file)
	write := bufio.NewWriter(file)
	_, err = write.WriteString(data)
	if err != nil {
		fmt.Println("写入文件缓存失败", err)
		return err
	}
	err = write.Flush()
	if err != nil {
		fmt.Println("写入文件失败", err)
		return err
	}
	return nil
}

func AppendToFile(filePath string, text string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("追加文件打开失败 %s!\n", file)
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("追加文件关闭失败", err)
		}
	}(file)
	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println("追加文件写入失败", err)
		return err
	}
	return nil
}

func GetPathOfSystemHostsPath() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("windir") + "\\system32\\drivers\\etc\\hosts"
	}
	return "/etc/hosts"
}
