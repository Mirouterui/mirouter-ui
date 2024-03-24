package database

import (
	"io"
	"strconv"

	"encoding/json"
	"fmt"
	"main/modules/config"
	"math"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// For database
type RouterHistory struct {
	gorm.Model
	Ip        string
	RouterNum int
	Cpu       float64
	Cpu_tp    int
	Mem       float64
	UpSpeed   float64
	DownSpeed float64
	UpTotal   float64
	DownTotal float64
	DeviceNum int
}
type DevicesHistory struct {
	gorm.Model
	Mac       string
	UpSpeed   float64
	DownSpeed float64
	UpTotal   float64
	DownTotal float64
}
type DeviceInfo struct {
	DevName          string `json:"devname"`
	Download         int64  `json:"download,string"`
	DownSpeed        int    `json:"downspeed,string"`
	Mac              string `json:"mac"`
	MaxDownloadSpeed int    `json:"maxdownloadspeed,string"`
	MaxUploadSpeed   int    `json:"maxuploadspeed,string"`
	Online           int    `json:"online,string"`
	Upload           int64  `json:"upload,string"`
	UpSpeed          int    `json:"upspeed,string"`
}

type Dev struct {
	Password   string `json:"password"`
	Key        string `json:"key"`
	IP         string `json:"ip"`
	RouterUnit bool   `json:"routerunit"`
}

// CheckDatabase checks the SQLite database file at the given path.
//
// Parameters:
// - databasePath: a string variable that holds the path to the SQLite database file.
//
// Returns: None.
func CheckDatabase(databasePath string) {
	// databasePath is a variable that holds the path to the SQLite database file
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	checkErr(err)

	// Check if the history table exists, if not, create it
	err = db.AutoMigrate(&RouterHistory{})
	checkErr(err)
	err = db.AutoMigrate(&DevicesHistory{})
	checkErr(err)
	// Perform CRUD operations on the history table using db.Create, db.First, db.Update, db.Delete methods
}

// Savetodb saves device statistics to the database.
//
// Parameters:
// - databasePath: the path to the database.
// - dev: an array of device configurations.
// - tokens: a map of token IDs to strings.
// - maxsaved: the maximum number of records to delete.
func Savetodb(databasePath string, dev []config.Dev, tokens map[int]string, maxsaved int) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	checkErr(err)
	for i, d := range dev {
		ip := d.IP
		routerNum := i
		cpu, cpu_tp, mem, upSpeed, downSpeed, upTotal, downTotal, deviceNum, devs := getRouterStats(i, tokens, ip)
		var count int64
		db.Model(&RouterHistory{}).Where("router_num = ?", routerNum).Count(&count)
		if count >= int64(maxsaved) {
			logrus.Debug("删除历史数据")
			db.Exec("DELETE FROM histories WHERE router_num = ? AND created_at = (SELECT MIN(created_at) FROM histories WHERE router_num = ? );", routerNum, routerNum)

		}
		db.Create(&RouterHistory{
			Ip:        ip,
			RouterNum: routerNum,
			Cpu:       cpu,
			Cpu_tp:    cpu_tp,
			Mem:       mem,
			UpSpeed:   upSpeed,
			DownSpeed: downSpeed,
			UpTotal:   upTotal,
			DownTotal: downTotal,
			DeviceNum: deviceNum,
		})
		for _, dev := range devs {
			devMap := dev.(map[string]interface{})

			data, err := json.Marshal(devMap)
			checkErr(err)

			var info DeviceInfo
			err = json.Unmarshal(data, &info)
			checkErr(err)
			mac := info.Mac
			upSpeed := float64(info.UpSpeed) / 1024 / 1024
			downSpeed := float64(info.DownSpeed) / 1024 / 1024
			upTotal := float64(info.Upload) / 1024 / 1024
			downTotal := float64(info.Download) / 1024 / 1024
			db.Create(&DevicesHistory{
				Mac:       mac,
				UpSpeed:   upSpeed,
				DownSpeed: downSpeed,
				UpTotal:   upTotal,
				DownTotal: downTotal,
			})
			db.Model(&DevicesHistory{}).Where("mac = ?", routerNum).Count(&count)
			if count >= int64(maxsaved) {
				logrus.Debug("删除历史数据")
				db.Exec("DELETE FROM histories WHERE mac = ? AND created_at = (SELECT MIN(created_at) FROM histories WHERE mac = ? );", mac, mac)

			}
		}
		db.Create(&RouterHistory{
			Ip:        ip,
			RouterNum: routerNum,
			Cpu:       cpu,
			Cpu_tp:    cpu_tp,
			Mem:       mem,
			UpSpeed:   upSpeed,
			DownSpeed: downSpeed,
			UpTotal:   upTotal,
			DownTotal: downTotal,
			DeviceNum: deviceNum,
		})

	}
}

func GetRouterHistory(databasePath string, routernum int, fixupfloat bool) []RouterHistory {

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	checkErr(err)
	var history []RouterHistory
	db.Where("router_num = ?", routernum).Find(&history)
	// 处理浮点数精度问题
	if fixupfloat {
		for i := range history {
			history[i].Cpu = round(history[i].Cpu, .5, 2)
			history[i].Mem = round(history[i].Mem, .5, 2)
			history[i].UpSpeed = round(history[i].UpSpeed, .5, 2)
			history[i].DownSpeed = round(history[i].DownSpeed, .5, 2)
			history[i].UpTotal = round(history[i].UpTotal, .5, 2)
			history[i].DownTotal = round(history[i].DownTotal, .5, 2)

		}

	}
	return history
}

func GetDeviceHistory(databasePath string, deviceMac string, fixupfloat bool) []DevicesHistory {

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	checkErr(err)
	var history []DevicesHistory
	db.Where("mac = ?", deviceMac).Find(&history)
	// 处理浮点数精度问题
	if fixupfloat {
		for i := range history {
			history[i].UpSpeed = round(history[i].UpSpeed, .5, 2)
			history[i].DownSpeed = round(history[i].DownSpeed, .5, 2)
			history[i].UpTotal = round(history[i].UpTotal, .5, 2)
			history[i].DownTotal = round(history[i].DownTotal, .5, 2)
		}

	}
	return history
}

func getRouterStats(routernum int, tokens map[int]string, ip string) (float64, int, float64, float64, float64, float64, float64, int, []interface{}) {
	if tokens[routernum] == "" {
		return 0, 0, 0, 0, 0, 0, 0, 0, []interface{}{}
	}
	url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/misystem/status", ip, tokens[routernum])
	resp, err := http.Get(url)
	checkErr(err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	upspeed, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["upspeed"].(string), 64)
	downspeed, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["downspeed"].(string), 64)
	uploadtotal, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["upload"].(string), 64)
	downloadtotal, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["download"].(string), 64)
	cpuload := result["cpu"].(map[string]interface{})["load"].(float64) * 100
	cpu_tp := int(result["temperature"].(float64))
	memusage := result["mem"].(map[string]interface{})["usage"].(float64) * 100
	devicenum_now := int(result["count"].(map[string]interface{})["online"].(float64))
	devs := result["dev"].([]interface{})

	return cpuload, cpu_tp, memusage, upspeed, downspeed, uploadtotal, downloadtotal, devicenum_now, devs
}

// func roundToOneDecimal(num float64) float64 {
// 	return math.Round(num*100) / 100
// }

func checkErr(err error) {
	if err != nil {
		logrus.Debug(err)
	}
}
func round(val float64, roundOn float64, places int) (newVal float64) {
	var rounder float64
	pow := math.Pow(10, float64(places))
	intermed := val * pow
	_, frac := math.Modf(intermed)

	if frac >= roundOn {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	newVal = rounder / pow
	return
}
