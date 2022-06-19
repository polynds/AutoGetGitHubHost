package main

import (
	"AutoGetGitHubHost/config"
	"fmt"
	"time"
)

var cfg *config.Config

func initConfig() {
	_cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return
	}
	//fmt.Printf("%+v\n", _cfg)
	cfg = _cfg
}

func main() {
	initConfig()

	if !cfg.Enabled {
		fmt.Println("配置未开启程序退出")
		return
	}

	ch := make(chan string, 2)
	//并发下载
	MultiDownload(cfg, ch)
	//hosts文件
	UpdateHosts(cfg, ch)

}

func UpdateHosts(config *config.Config, ch chan string) {
	timeout := time.After(900 * time.Second)
	for idx := 0; idx < len(config.Hosts); idx++ {
		select {
		case filePath := <-ch:
			nt := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s]Finish download %s\n", nt, filePath)
			err := updateOsHosts(filePath)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		case <-timeout:
			fmt.Println("Timeout...")
			break
		}
	}
}
