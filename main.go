package main

import (
	"AutoGetGitHubHost/config"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var cfg *config.Config

func initConfig() {
	_cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return
	}
	cfg = _cfg
}

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			fmt.Println("close...")
			return
		case <-t.ticker.C:
			t.wg.Add(1)
			fmt.Println("doUpdateHosts...")
			go func() {
				defer t.wg.Done()
				doUpdateHosts()
			}()
		}
	}
}

func (t *Task) Stop() {
	fmt.Println("got close signal")
	close(t.closed)
	//在这里会等待所有的协程都退出
	t.wg.Wait()
	fmt.Println("all goroutine　done")
}

func main() {
	initConfig()

	if !cfg.Enabled {
		fmt.Println("配置未开启程序退出")
		return
	}

	task := &Task{
		closed: make(chan struct{}),
		ticker: time.NewTicker(time.Hour * 1),
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go task.Run()

	go doUpdateHosts()

	select {
	case sig := <-c:

		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.Stop()
	}
}

func doUpdateHosts() {
	var ch = make(chan string, 2)
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
