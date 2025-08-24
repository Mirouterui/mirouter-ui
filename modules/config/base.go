package config

import (
	"errors"
	"flag"
	"io"
	. "main/modules/download"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Mirouterui/mirouter-ui/modules/download"
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
	basedirectory     string
	databasepath      string
	historyEnable     bool
	dev               []Dev
	maxsaved          int
	flushTokenTime    int
	sampletime        int
	netdata_routernum int
)

type Dev struct {
	Password   string `mapstructure:"password"`
	Key        string `mapstructure:"key"`
	IP         string `mapstructure:"ip"`
	RouterUnit bool   `mapstructure:"routerunit"`
	IsLocal    bool   `mapstructure:"islocal"`
}

type History struct {
	Enable     bool `mapstructure:"enable"`
	MaxDeleted int  `mapstructure:"maxsaved"`
	Sampletime int  `mapstructure:"sampletime"`
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

func GetConfigInfo() (dev []Dev, debug bool, port int, tiny bool, basedirectory string, flushTokenTime int, databasepath string, maxsaved int, historyEnable bool, sampletime int, netdata_routernum int) {
	flag.StringVar(&configPath, "config", "", "配置文件路径")
	flag.StringVar(&basedirectory, "basedirectory", "", "基础目录路径")
	flag.StringVar(&databasepath, "databasepath", "", "数据库路径")
	flag.Parse()
	appPath, err := os.Executable()
	checkErr(err)
	if configPath == "" {
		configPath = filepath.Join(filepath.Dir(appPath), "config.json")
	}
	if databasepath == "" {
		databasepath = filepath.Join(filepath.Dir(appPath), "database.db")
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
	maxsaved = config.History.MaxDeleted
	historyEnable = config.History.Enable
	sampletime = config.History.Sampletime
	flushTokenTime = config.FlushTokenTime
	netdata_routernum = config.Netdata_routernum
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
	return dev, debug, port, tiny, basedirectory, flushTokenTime, databasepath, maxsaved, historyEnable, sampletime, netdata_routernum
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
