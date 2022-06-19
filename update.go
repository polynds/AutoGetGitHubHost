package main

import (
	"AutoGetGitHubHost/utils"
	"errors"
	"fmt"
	"regexp"
	"time"
)

const (
	IPReg      = `[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}`
	UrlReg     = `[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+`
	Header     = "#AutoGetGitHubHostStart"
	Footer     = "#AutoGetGitHubHostEnd"
	ReplaceReg = Header + `([\w\W]+)` + Footer
)

func updateOsHosts(filePath string) error {
	if utils.IsExist(filePath) {
		//读取全部
		text := utils.ReadAll(filePath)
		hostsContent := formatContent(string(text))
		fmt.Println(hostsContent)
		replaceContent(hostsContent)
		return nil
	}
	return errors.New("下载的hosts文件未找到,程序退出")
}

func formatContent(text string) string {
	//匹配出所有的ip和域名的map
	reg := regexp.MustCompile(IPReg)
	ips := reg.FindAllString(text, -1)
	fmt.Printf("%q\n", ips)

	reg = regexp.MustCompile(UrlReg)
	urls := reg.FindAllString(text, -1)
	fmt.Printf("%q\n", urls)
	fmt.Println(len(ips), len(urls))
	//然后拼接成前后标志的字符串
	content := ""
	for i := 0; i < len(ips); i++ {
		ip := ips[i]
		url := urls[i]
		if ip == "" || url == "" {
			continue
		}
		content += ips[i] + "  " + urls[i] + "\n"
	}
	content += time.Now().Format("2006-01-02 15:04:05") + "\n"
	//替换或者追加到hosts文件末尾
	return fmt.Sprintf("%s\n%s%s\n", Header, content, Footer)
}

func replaceContent(hostsContent string) {
	systemHostsPath := utils.GetPathOfSystemHostsPath()

	systemHostsText := utils.ReadAll(systemHostsPath)

	reg := regexp.MustCompile(ReplaceReg)
	text := reg.FindAllString(string(systemHostsText), -1)
	if len(text) < 1 {
		//还没有追加到文件末尾
		err := utils.AppendToFile(systemHostsPath, hostsContent)
		if err != nil {
			fmt.Println("写入hosts文件内容失败")
			return
		}
	} else {
		//替换
		afterText := reg.ReplaceAllString(hostsContent, `${1}`)
		err := utils.WriteAll(systemHostsPath, afterText)
		if err != nil {
			fmt.Println("替换hosts文件内容失败")
			return
		}
	}
	fmt.Println(string(systemHostsText), text)
}
