package main

import (
	"AutoGetGitHubHost/config"
	"AutoGetGitHubHost/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const DownloadDir = "./hosts/"

func getFilePath(url string) string {
	fmt.Println(getFileName(url))
	return fmt.Sprintf("%s%s", DownloadDir, getFileName(url))
}

func getFileName(url string) string {
	return fmt.Sprintf("%s.txt", utils.Str2md5(url))
}
func MultiDownload(config *config.Config, ch chan string) {
	hosts := config.Hosts
	var wg sync.WaitGroup
	wg.Add(len(hosts))
	for _, url := range hosts {
		go func(url string) {
			defer wg.Done()
			ch <- download(url)
		}(url)
	}
	wg.Wait()
}

func download(url string) string {

	client := http.Client{Timeout: 900 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(response.Body)

	filePath := getFilePath(url)
	fmt.Println(url, filePath)
	utils.SaveFile(filePath, response.Body)

	return filePath
}
