package main

import (
	"encoding/json"
	_ "flag"
	"fmt"
	"io"
	"main/modules/config"
	"main/modules/database"
	"main/modules/download"
	login "main/modules/login"
	"main/modules/netdata"
	"main/modules/tp"
	"net/http"

	// _ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
)

var (
	password          string
	key               string
	token             string
	tokens            map[int]string
	debug             bool
	port              int
	routerName        string
	routerNames       map[int]string
	hardware          string
	hardwares         map[int]string
	routerunits       map[int]bool
	tiny              bool
	routerunit        bool
	dev               []config.Dev
	cpu_cmd           *exec.Cmd
	w24g_cmd          *exec.Cmd
	w5g_cmd           *exec.Cmd
	configPath        string
	workdirectory     string
	Version           string
	databasepath      string
	flushTokenTime    int
	maxsaved          int
	historyEnable     bool
	sampletime        int
	netdata_routernum int
)

type Config struct {
	Dev          []config.Dev `json:"dev"`
	Debug        bool         `json:"debug"`
	Port         int          `json:"port"`
	Tiny         bool         `json:"tiny"`
	Databasepath string       `json:"databasepath"`
}

func init() {
	dev, debug, port, tiny, workdirectory, flushTokenTime, databasepath, maxsaved, historyEnable, sampletime, netdata_routernum = config.GetConfigInfo()
	tokens = make(map[int]string)
	routerNames = make(map[int]string)
	hardwares = make(map[int]string)
	routerunits = make(map[int]bool)
	// go func() {
	// 	logrus.Println(http.ListenAndServe(":6060", nil))
	// }()
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
	type History struct {
		Enable       bool   `json:"enable"`
		MaxDeleted   int    `json:"maxsaved"`
		Databasepath string `json:"databasepath"`
		Sampletime   int    `json:"sampletime"`
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
	history := History{}
	history.Enable = historyEnable
	history.MaxDeleted = maxsaved
	history.Databasepath = databasepath
	history.Sampletime = sampletime
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":  0,
		"tiny":  tiny,
		"port":  port,
		"debug": debug,
		// "token":      token,
		"dev":            devsNoPassword,
		"history":        history,
		"flushTokenTime": flushTokenTime,
		"ver":            Version,
	})
}

func gettoken(dev []config.Dev) {
	for i, d := range dev {
		token, routerName, hardware := login.GetToken(d.Password, d.Key, d.IP)
		tokens[i] = token
		routerNames[i] = routerName
		hardwares[i] = hardware
		routerunits[i] = d.RouterUnit
		logrus.Debug(hardwares[i])
	}
}
func main() {
	starttime := int(time.Now().Unix())
	logrus.Info("当前后端版本为：" + Version)
	e := echo.New()
	c := cron.New()
	e.Use(middleware.Recover())
	// 输出访问日志
	if debug {
		e.Use(middleware.Logger())
	}
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		c.Response().Header().Set("Access-Control-Allow-Private-Network", "true")
	// 		return next(c)
	// 	}
	// })

	e.Use(middleware.CORS())

	e.GET("/:routernum/api/:apipath", func(c echo.Context) error {
		routernum, err := strconv.Atoi(c.Param("routernum"))
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 1100, "msg": "参数错误"})
		}
		apipath := c.Param("apipath")
		ip := dev[routernum].IP

		switch apipath {

		case "xqsystem/router_name":
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 0, "routerName": routerNames[routernum]})

		case "misystem/status", "misystem/devicelist", "xqsystem/internet_connect", "xqsystem/fac_info", "misystem/messages", "xqsystem/upnp":
			url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/%s", ip, tokens[routernum], apipath)
			resp, err := http.Get(url)
			if err != nil {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"code": 1101,
					"msg":  "小米路由器的api调用出错，请检查配置或路由器状态",
				})
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			if routerunits[routernum] && apipath == "misystem/status" {
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

	e.GET("/:routernum/_api/gettemperature", func(c echo.Context) error {
		routernum, err := strconv.Atoi(c.Param("routernum"))
		logrus.Debug(tokens)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 1100, "msg": "参数错误"})
		}
		status, cpu_tp, fanspeed, w24g_tp, w5g_tp := tp.GetTemperature(c, routernum, hardwares[routernum])
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

	e.GET("/api/v1/data", func(c echo.Context) error {
		chart := c.QueryParam("chart")
		dimensions := c.QueryParam("dimensions")

		ip := dev[netdata_routernum].IP
		token := tokens[netdata_routernum]
		cpuLoad, memAvailable, _, _, upSpeed, downSpeed, temperature, deviceOnline, _, _ := netdata.ProcessData(ip, token)

		switch chart {

		case "system.cpu":
			if routerunits[netdata_routernum] {
				cpuLoad = int(GetCpuPercent() * 100)
			}
			data := netdata.GenerateArray("system.cpu", cpuLoad, starttime, "system.cpu", "system.cpu")
			return c.JSON(http.StatusOK, data)
		case "mem.available":
			data := netdata.GenerateArray("mem.available", memAvailable, starttime, "avail", "MemAvailable")
			return c.JSON(http.StatusOK, data)
		case "device.online":
			data := netdata.GenerateArray("device.online", deviceOnline, starttime, "online", "online")
			return c.JSON(http.StatusOK, data)
		case "net.eth0":
			if dimensions == "received" {
				data := netdata.GenerateArray("net.eth0", downSpeed, starttime, "received", "received")
				return c.JSON(http.StatusOK, data)
			}
			if dimensions == "sent" {
				data := netdata.GenerateArray("net.eth0", -upSpeed, starttime, "sent", "sent")
				return c.JSON(http.StatusOK, data)
			}
			return c.String(http.StatusOK, "缺失参数")
		case "sensors.temp_thermal_zone0_thermal_thermal_zone0":
			data := netdata.GenerateArray("sensors.temp_thermal_zone0_thermal_thermal_zone0", temperature, starttime, "temperature", "temperature")
			return c.JSON(http.StatusOK, data)
		default:
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1102,
				"msg":  "该图表数据不支持",
			})
		}
	})
	// 没用
	// e.GET("/api/v1/charts", func(c echo.Context) error {
	// 	ip := dev[netdata_routernum].IP
	// 	token := tokens[netdata_routernum]
	// 	cpuLoad, memAvailable, memTotal, memUsage, upSpeed, downSpeed, temperature, deviceonline := netdata.ProcessData(ip, token)
	// 	cpuLoadData := netdata.GenerateDataForAllMetrics("system.cpu", "cpu", "percentage", cpuLoad, "used")
	// 	memAvailableData := netdata.GenerateDataForAllMetrics("mem.available", "mem", "bytes", memAvailable, "used")
	// 	memTotalData := netdata.GenerateDataForAllMetrics("mem.total", "mem", "bytes", memTotal, "used")
	// 	memUsageData := netdata.GenerateDataForAllMetrics("mem.used", "mem", "percentage", memUsage, "used")
	// 	upSpeedData := netdata.GenerateDataForAllMetrics("net.eth0.receivedspeed", "net", "bytes", upSpeed, "received")
	// 	downSpeedData := netdata.GenerateDataForAllMetrics("net.eth0.sentspeed", "net", "bytes", downSpeed, "sent")
	// 	temperatureData := netdata.GenerateDataForAllMetrics("sensors.temp_thermal_zone0_thermal_thermal_zone0", "sensors", "celsius", temperature, "temperature")
	// 	deviceonlineData := netdata.GenerateDataForAllMetrics("device.online", "device", "count", deviceonline, "online")
	// 	charts := map[string]interface{}{
	// 		"system.cpu":             cpuLoadData,
	// 		"mem.available":          memAvailableData,
	// 		"mem.total":              memTotalData,
	// 		"mem.used":               memUsageData,
	// 		"net.eth0.receivedspeed": upSpeedData,
	// 		"net.eth0.sentspeed":     downSpeedData,
	// 		"device.online":          deviceonlineData,
	// 		"sensors.temp_thermal_zone0_thermal_thermal_zone0": temperatureData,
	// 	}
	// 	data := map[string]interface{}{
	// 		"hostname":        routerNames[netdata_routernum],
	// 		"version":         "v1.29.3",
	// 		"release_channel": "stable",
	// 		"os":              "linux",
	// 		"timezone":        "Asia/Shanghai",
	// 		"update_every":    1,
	// 		"history":         3996,
	// 		"memory_mode":     "dbengine",
	// 		"custom_info":     "",
	// 		"charts":          charts,
	// 	}

	// 	return c.JSON(http.StatusOK, data)

	// })

	// 应付HA用
	e.GET("/api/v1/allmetrics?format=json&help=no&types=no&timestamps=yes&names=yes&data=average", func(c echo.Context) error {
		ip := dev[netdata_routernum].IP
		token := tokens[netdata_routernum]
		cpuLoad, memAvailable, memTotal, memUsage, upSpeed, downSpeed, temperature, deviceonline, uploadtotal, downloadtotal := netdata.ProcessData(ip, token)
		cpuLoadData := netdata.GenerateDataForAllMetrics("system.cpu", "cpu", "percentage", cpuLoad, "used")
		memAvailableData := netdata.GenerateDataForAllMetrics("mem.available", "mem", "bytes", memAvailable, "used")
		memTotalData := netdata.GenerateDataForAllMetrics("mem.total", "mem", "bytes", memTotal, "used")
		memUsageData := netdata.GenerateDataForAllMetrics("mem.used", "mem", "percentage", memUsage, "used")
		upSpeedData := netdata.GenerateDataForAllMetrics("net.eth0.sentspeed", "net", "bytes", upSpeed, "sent")
		downSpeedData := netdata.GenerateDataForAllMetrics("net.eth0.receivedspeed", "net", "bytes", downSpeed, "received")
		temperatureData := netdata.GenerateDataForAllMetrics("sensors.temp_thermal_zone0_thermal_thermal_zone0", "sensors", "celsius", temperature, "temperature")
		deviceonlineData := netdata.GenerateDataForAllMetrics("device.online", "device", "count", deviceonline, "online")
		uploadtotalData := netdata.GenerateDataForAllMetrics("net.eth0.sent", "net", "bytes", uploadtotal, "total")
		downloadtotalData := netdata.GenerateDataForAllMetrics("net.eth0.received", "net", "bytes", downloadtotal, "total")
		data := map[string]interface{}{
			"system.cpu":             cpuLoadData,
			"mem.available":          memAvailableData,
			"mem.total":              memTotalData,
			"mem.used":               memUsageData,
			"net.eth0.receivedspeed": downSpeedData,
			"net.eth0.sentspeed":     upSpeedData,
			"device.online":          deviceonlineData,
			"net.eth0.sent":          uploadtotalData,
			"net.eth0.received":      downloadtotalData,
			"sensors.temp_thermal_zone0_thermal_thermal_zone0": temperatureData,
		}

		return c.JSON(http.StatusOK, data)

	})
	// e.GET("/api/v1/alarms?all&format=json", func(c echo.Context) error {
	// 	time := int(time.Now().Unix())
	// 	var value int
	// 	var status string
	// 	if login.CheckRouterAvailability(dev[netdata_routernum].IP) {
	// 		value = 1
	// 		status = "CLEAR"
	// 	} else {
	// 		value = 0
	// 		status = "CRITICAL"
	// 	}
	// 	alarm := map[string]interface{}{
	// 		"id":                    1,
	// 		"name":                  "router_offline",
	// 		"chart":                 "router.status",
	// 		"family":                "status",
	// 		"active":                true,
	// 		"disabled":              false,
	// 		"silenced":              false,
	// 		"exec":                  "/usr/lib/netdata/plugins.d/alarm-notify.sh",
	// 		"recipient":             "sysadmin",
	// 		"source":                "10@/usr/lib/netdata/conf.d/health.d/router_offline.conf",
	// 		"units":                 "status",
	// 		"info":                  "the status of the router (offline = 0, online = 1)",
	// 		"status":                status,
	// 		"last_status_change":    1704026010,
	// 		"last_updated":          time,
	// 		"next_update":           time + 10,
	// 		"update_every":          10,
	// 		"delay_up_duration":     0,
	// 		"delay_down_duration":   300,
	// 		"delay_max_duration":    3600,
	// 		"delay_multiplier":      1.5,
	// 		"delay":                 0,
	// 		"delay_up_to_timestamp": 1704026010,
	// 		"warn_repeat_every":     "0",
	// 		"crit_repeat_every":     "0",
	// 		"value_string":          "1",
	// 		"last_repeat":           "0",
	// 		"calc":                  "${status}",
	// 		"calc_parsed":           "${status}",
	// 		"warn":                  "$this == 0",
	// 		"warn_parsed":           "${this} == 0",
	// 		"crit":                  "$this == 0",
	// 		"crit_parsed":           "${this} == 0",
	// 		"green":                 nil,
	// 		"red":                   nil,
	// 		"value":                 value,
	// 	}

	// 	data := map[string]interface{}{
	// 		"hostname":                   routerNames[netdata_routernum],
	// 		"latest_alarm_log_unique_id": 1703857080,
	// 		"status":                     true,
	// 		"now":                        time,
	// 		"alarms": map[string]interface{}{
	// 			"router_offline": alarm,
	// 		},
	// 	}
	// 	return c.JSON(http.StatusOK, data)
	// })
	e.GET("/_api/getconfig", getconfig)

	e.GET("/_api/getrouterhistory", func(c echo.Context) error {
		routernum, err := strconv.Atoi(c.QueryParam("routernum"))
		fixupfloat := c.QueryParam("fixupfloat")
		if fixupfloat == "" {
			fixupfloat = "false"
		}
		fixupfloat_bool, err1 := strconv.ParseBool(fixupfloat)
		if err != nil || err1 != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 1100, "msg": "参数错误"})
		}
		if !historyEnable {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1101,
				"msg":  "历史数据未开启",
			})
		}
		history := database.GetRouterHistory(databasepath, routernum, fixupfloat_bool)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"history": history,
		})
	})

	e.GET("/_api/getdevicehistory", func(c echo.Context) error {
		deviceMac := c.QueryParam("devicemac")
		fixupfloat := c.QueryParam("fixupfloat")
		if fixupfloat == "" {
			fixupfloat = "false"
		}
		fixupfloat_bool, err := strconv.ParseBool(fixupfloat)

		if deviceMac == "" || len(deviceMac) != 17 || err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 1100, "msg": "参数错误"})
		}
		if !historyEnable {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1101,
				"msg":  "历史数据未开启",
			})
		}
		history := database.GetDeviceHistory(databasepath, deviceMac, fixupfloat_bool)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"history": history,
		})
	})
	e.GET("/_api/flushstatic", func(c echo.Context) error {
		err := download.DownloadStatic(workdirectory, true, true)
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

	e.GET("/_api/refresh", func(c echo.Context) error {
		gettoken(dev)
		logrus.Debugln("执行完成")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "执行完成",
		})
	})
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

	// var contentHandler = echo.WrapHandler(http.FileServer(http.FS(static)))
	// var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/static/$1"})

	// e.GET("/*", contentHandler, contentRewrite)
	if !tiny {
		directory := "static"
		if workdirectory != "" {
			directory = filepath.Join(workdirectory, "static")
		}
		logrus.Debug("静态资源目录为:" + directory)
		e.Static("/", directory)
	} else if tiny {
		e.GET("/*", func(c echo.Context) error {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"code": 404, "msg": "已开启tiny模式"})
		})
	}
	gettoken(dev)

	database.CheckDatabase(databasepath)
	c.AddFunc("@every "+strconv.Itoa(flushTokenTime)+"s", func() { gettoken(dev) })

	if historyEnable {
		c.AddFunc("@every "+strconv.Itoa(sampletime)+"s", func() { database.Savetodb(databasepath, dev, tokens, maxsaved) })
	}
	c.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		e.Close()
	}()

	e.Start(":" + fmt.Sprint(port))
}
