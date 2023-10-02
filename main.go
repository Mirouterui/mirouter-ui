package main

import (
	"encoding/json"
	_ "flag"
	"fmt"
	"io"
	"main/modules/config"
	"main/modules/download"
	login "main/modules/login"
	"main/modules/tp"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

var (
	password      string
	key           string
	ip            string
	token         string
	tokens        map[int]string
	debug         bool
	port          int
	routerName    string
	routerNames   map[int]string
	hardware      string
	hardwares     map[int]string
	tiny          bool
	routerunit    bool
	dev           []config.Dev
	cpu_cmd       *exec.Cmd
	w24g_cmd      *exec.Cmd
	w5g_cmd       *exec.Cmd
	configPath    string
	basedirectory string
	Version       string
)

type Config struct {
	Dev   []config.Dev `json:"dev"`
	Debug bool         `json:"debug"`
	Port  int          `json:"port"`
	Tiny  bool         `json:"tiny"`
}

func init() {
	dev, debug, port, tiny, basedirectory = config.Getconfig()
	tokens = make(map[int]string)
	routerNames = make(map[int]string)
	hardwares = make(map[int]string)
}
func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0] / 100
}

func getconfig(c echo.Context) error {
	type DevNoPassword struct {
		Key        string `json:"key"`
		IP         string `json:"ip"`
		RouterUnit bool   `json:"routerunit"`
	}

	devsNoPassword := []DevNoPassword{}
	for _, d := range dev {
		devNoPassword := DevNoPassword{
			Key:        d.Key,
			IP:         d.IP,
			RouterUnit: d.RouterUnit,
		}
		devsNoPassword = append(devsNoPassword, devNoPassword)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":  0,
		"tiny":  tiny,
		"port":  port,
		"debug": debug,
		// "token":      token,
		"dev": devsNoPassword,
		"ver": Version,
	})
}

func gettoken(dev []config.Dev) {
	for i, d := range dev {
		token, routerName, hardware := login.GetToken(d.Password, d.Key, d.IP)
		tokens[i] = token
		routerNames[i] = routerName
		hardwares[i] = hardware
		logrus.Debug(hardwares[i])
	}
}
func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		c.Response().Header().Set("Access-Control-Allow-Private-Network", "true")
	// 		return next(c)
	// 	}
	// })

	e.Use(middleware.CORS())

	e.GET("/:devnum/api/:apipath", func(c echo.Context) error {
		devnum, err := strconv.Atoi(c.Param("devnum"))
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 1100, "msg": "参数错误"})
		}
		apipath := c.Param("apipath")
		ip = dev[devnum].IP

		switch apipath {

		case "xqsystem/router_name":
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 0, "routerName": routerNames[devnum]})

		case "misystem/status", "misystem/devicelist", "xqsystem/internet_connect", "xqsystem/fac_info", "misystem/messages":
			url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/%s", ip, tokens[devnum], apipath)
			resp, err := http.Get(url)
			if err != nil {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"code": 1101,
					"msg":  "MiRouter的api调用出错，请检查配置或路由器状态",
				})
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			if routerunit || apipath == "misystem/status" {
				cpuPercent := GetCpuPercent()
				if cpu, ok := result["cpu"].(map[string]interface{}); ok {
					cpu["load"] = cpuPercent
				}
			}
			return c.JSON(http.StatusOK, result)

		default:
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1102,
				"msg":  "该api不支持免密调用",
			})
		}
	})

	e.GET("/:devnum/_api/gettemperature", func(c echo.Context) error {
		devnum, err := strconv.Atoi(c.Param("devnum"))
		logrus.Debug(tokens)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 1100, "msg": "参数错误"})
		}
		status, cpu_tp, fanspeed, w24g_tp, w5g_tp := tp.GetTemperature(c, devnum, hardwares[devnum])
		if status {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code":     0,
				"cpu":      cpu_tp,
				"fanspeed": fanspeed,
				"w24g":     w24g_tp,
				"w5g":      w5g_tp,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 1103,
			"msg":  "不支持该设备",
		})
	})

	e.GET("/_api/getconfig", getconfig)
	e.GET("/_api/quit", func(c echo.Context) error {
		go func() {
			time.Sleep(1 * time.Second)
			defer os.Exit(0)
		}()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "正在关闭",
		})
	})
	e.GET("/_api/flushstatic", func(c echo.Context) error {
		err := download.DownloadStatic(basedirectory, true)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1101,
				"msg":  err,
			})
		}
		logrus.Debugln("执行完成")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "执行完成",
		})
	})

	// var contentHandler = echo.WrapHandler(http.FileServer(http.FS(static)))
	// var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/static/$1"})

	// e.GET("/*", contentHandler, contentRewrite)
	if tiny == false {
		directory := "static"
		if basedirectory != "" {
			directory = filepath.Join(basedirectory, "static")
		}
		logrus.Debug("静态资源目录为:" + directory)
		e.Static("/", directory)
	}
	gettoken(dev)

	go func() {
		for range time.Tick(30 * time.Minute) {
			gettoken(dev)
		}
	}()
	e.Start(":" + fmt.Sprint(port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		e.Close()
	}()
}
