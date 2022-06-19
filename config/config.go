package config

import (
	"AutoGetGitHubHost/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	LogFilePath        string = "./error.log"
	JsonConfigFilePath string = "./config/config.json"
	DefaultConfig      string = `{
	"enabled": true,
	"hosts": [
		"https://cdn.jsdelivr.net/gh/521xueweihan/GitHub520@main/hosts"
	]
}`
)

var fileLocker sync.Mutex

type Config struct {
	Enabled bool     `json:"enabled"`
	Hosts   []string `json:"hosts"`
}

func NewConfig() *Config {
	return &Config{}
}

func init() {
	initLogFile()
	initConfigFile()
}

func initConfigFile() {
	if utils.CheckAndCreateDir(JsonConfigFilePath) {
		fmt.Println("json配置文件夹创建失败")
		return
	}
	if !utils.IsExist(JsonConfigFilePath) {
		ioReader := ioutil.NopCloser(strings.NewReader(DefaultConfig))
		utils.SaveFile(JsonConfigFilePath, ioReader)
	}
}

func initLogFile() {
	logfile, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetPrefix("Trace: ")
	log.SetOutput(logfile)
}

var _Config = NewConfig()

func InitConfig() (*Config, error) {

	fileLocker.Lock()
	file, err := ioutil.ReadFile(JsonConfigFilePath)
	if err != nil {
		log.Println(err)
		return _Config, err
	}

	err = json.Unmarshal(file, _Config)
	if err != nil {
		log.Println(err)
		return _Config, err
	}

	return _Config, nil
}
