package config

import (
	"encoding/json"
	"flag"
	. "main/modules/download"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	password          string
	key               string
	ip                string
	debug             bool
	port              int
	tiny              bool
	routerunit        bool
	configPath        string
	workdirectory     string
	databasepath      string
	historyEnable     bool
	dev               []Dev
	maxsaved          int
	flushTokenTime    int
	sampletime        int
	netdata_routernum int
	autocheckupdate   string
)

type Dev struct {
	Password   string `json:"password"`
	Key        string `json:"key"`
	IP         string `json:"ip"`
	RouterUnit bool   `json:"routerunit"`
}
type History struct {
	Enable     bool `json:"enable"`
	MaxDeleted int  `json:"maxsaved"`
	Sampletime int  `json:"sampletime"`
}
type Config struct {
	Dev               []Dev   `json:"dev"`
	History           History `json:"history"`
	Debug             bool    `json:"debug"`
	Port              int     `json:"port"`
	Tiny              bool    `json:"tiny"`
	FlushTokenTime    int     `json:"flushTokenTime"`
	Netdata_routernum int     `json:"netdata_routernum"`
}

func GetConfigInfo() (dev []Dev, debug bool, port int, tiny bool, workdirectory string, flushTokenTime int, databasepath string, maxsaved int, historyEnable bool, sampletime int, netdata_routernum int) {
	appPath, _ := os.Executable()

	flag.StringVar(&configPath, "config", filepath.Join(filepath.Dir(appPath), "config.json"), "配置文件路径")
	flag.StringVar(&workdirectory, "workdirectory", "", "工作目录路径")
	flag.StringVar(&databasepath, "databasepath", filepath.Join(filepath.Dir(appPath), "database.db"), "数据库路径")
	flag.StringVar(&autocheckupdate, "autocheckupdate", "true", "自动检查更新")
	flag.Parse()

	autocheckupdatebool, _ := strconv.ParseBool(autocheckupdate)

	logrus.Info("配置文件路径为:" + configPath)
	data, err := os.ReadFile(configPath)

	if err != nil {
		logrus.Info("未找到配置文件，正在从程序内部导出")
		// 使用你的结构体创建一个默认的配置实例
		config := Config{
			Dev: []Dev{
				{
					Password:   "",
					Key:        "a2ffa5c9be07488bbb04a3a47d3c5f6a",
					IP:         "192.168.31.1",
					RouterUnit: false,
				},
			},
			History: History{
				Enable:     false,
				MaxDeleted: 3000,
				Sampletime: 86400,
			},
			Debug:             true,
			Port:              6789,
			Tiny:              false,
			FlushTokenTime:    1800,
			Netdata_routernum: 0,
		}
		configContent, err := json.MarshalIndent(config, "", "  ")
		checkErr(err)
		err = os.WriteFile(configPath, configContent, 0644)
		checkErr(err)
		logrus.Info("配置文件导出完成，请修改配置文件")
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
	maxsaved = config.History.MaxDeleted
	historyEnable = config.History.Enable
	sampletime = config.History.Sampletime
	flushTokenTime = config.FlushTokenTime
	netdata_routernum = config.Netdata_routernum
	// logrus.Info(password)
	// logrus.Info(key)
	if !tiny {
		DownloadStatic(workdirectory, false, autocheckupdatebool)
	}
	if debug {
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
	return dev, debug, port, tiny, workdirectory, flushTokenTime, databasepath, maxsaved, historyEnable, sampletime, netdata_routernum
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
