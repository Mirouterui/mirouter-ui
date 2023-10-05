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
	password       string
	key            string
	ip             string
	debug          bool
	port           int
	tiny           bool
	routerunit     bool
	configPath     string
	basedirectory  string
	databasepath   string
	historyEnable  bool
	Version        string
	dev            []Dev
	maxdeleted     int64
	flushTokenTime int64
	sampletime     int64
)

type Dev struct {
	Password   string `json:"password"`
	Key        string `json:"key"`
	IP         string `json:"ip"`
	RouterUnit bool   `json:"routerunit"`
}
type History struct {
	Enable       bool   `json:"enable"`
	MaxDeleted   int64  `json:"maxdeleted"`
	Databasepath string `json:"databasepath"`
	Sampletime   int64  `json:"sampletime"`
}
type Config struct {
	Dev            []Dev   `json:"dev"`
	History        History `json:"history"`
	Debug          bool    `json:"debug"`
	Port           int     `json:"port"`
	Tiny           bool    `json:"tiny"`
	FlushTokenTime int64   `json:"flushTokenTime"`
}

func GetConfigInfo() (dev []Dev, debug bool, port int, tiny bool, basedirectory string, databasepath string, flushTokenTime int64, maxdeleted int64, historyEnable bool, sampletime int64) {
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

	logrus.Info("配置文件路径为:" + configPath)
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
	dev = config.Dev
	debug = config.Debug
	port = config.Port
	tiny = config.Tiny
	databasepath = config.History.Databasepath
	maxdeleted = config.History.MaxDeleted
	historyEnable = config.History.Enable
	sampletime = config.History.Sampletime
	flushTokenTime = config.FlushTokenTime
	// logrus.Info(password)
	// logrus.Info(key)
	if tiny == false {
		DownloadStatic(basedirectory, false)
	}
	if debug == true {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	numDevs := len(dev)
	if numDevs == 0 {
		logrus.Info("未填写路由器信息，请检查配置文件")
		logrus.Info("5秒后退出程序")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	return dev, debug, port, tiny, basedirectory, databasepath, flushTokenTime, maxdeleted, historyEnable, sampletime
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
