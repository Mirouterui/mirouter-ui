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
	dev           []Dev
)

type Dev struct {
	Password   string `json:"password"`
	Key        string `json:"key"`
	IP         string `json:"ip"`
	RouterUnit bool   `json:"routerunit"`
}
type Config struct {
	Dev   []Dev `json:"dev"`
	Debug bool  `json:"debug"`
	Port  int   `json:"port"`
	Tiny  bool  `json:"tiny"`
}

func Getconfig() (dev []Dev, debug bool, port int, tiny bool, basedirectory string) {
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
	return dev, debug, port, tiny, basedirectory
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
