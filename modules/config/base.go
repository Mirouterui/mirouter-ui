package config

import (
	"encoding/json"
	"flag"
	"io"
	. "main/modules/download"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	password      string
	key           string
	ip            string
	debug         bool
	port          int
	tiny          bool
	routerunit    bool
	configPath    string
	basedirectory string
	Version       string
)

type Config struct {
	Password   string `json:"password"`
	Key        string `json:"key"`
	Ip         string `json:"ip"`
	Debug      bool   `json:"debug"`
	Port       int    `json:"port"`
	Tiny       bool   `json:"tiny"`
	Routerunit bool   `json:"routerunit"`
}

func init() {
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.StringVar(&basedirectory, "basedirectory", "", "基础目录路径")
	flag.Parse()
	if configPath == "" {
		appPath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		configPath = filepath.Join(filepath.Dir(appPath), "config.json")
	}

	// logrus.Info(configPath)
	logrus.Debug("配置文件路径为:" + configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Info("未找到配置文件，正在下载")
		resp, err := http.Get("https://mrui-api.hzchu.top/downloadconfig")
		checkErr(err)
		defer resp.Body.Close()
		out, err := os.Create(configPath)
		checkErr(err)
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		checkErr(err)
		logrus.Info("下载配置文件完成，请修改配置文件")
		logrus.Info("5秒后退出程序")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		logrus.Info("配置文件存在错误")
	}
	password = config.Password
	key = config.Key
	ip = config.Ip
	debug = config.Debug
	port = config.Port
	tiny = config.Tiny
	routerunit = config.Routerunit
	// logrus.Info(password)
	// logrus.Info(key)
	// logrus.Info(iv)
	if tiny == false {
		CheckAndDownloadStatic(basedirectory)
	}
	if debug == true {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
