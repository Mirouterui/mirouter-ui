package database

import (
	"io"
	"strconv"

	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"

	// "math"
	"net/http"
	"time"

	"github.com/Mirouterui/mirouter-ui/modules/config"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// For database
type RouterHistory struct {
	gorm.Model
	Ip        string
	RouterNum int
	Cpu       int
	Cpu_tp    int
	Mem       int
	UpSpeed   int
	DownSpeed int
	UpTotal   int
	DownTotal int
	DeviceNum int
}
type DevicesHistory struct {
	gorm.Model
	Mac       string
	UpSpeed   int
	DownSpeed int
	UpTotal   int
	DownTotal int
}
type DeviceInfo struct {
	DevName          string `json:"devname"`
	DownloadTotal    int    `json:"download,string"`
	DownSpeed        int    `json:"downspeed,string"`
	Mac              string `json:"mac"`
	MaxDownloadSpeed int    `json:"maxdownloadspeed,string"`
	MaxUploadSpeed   int    `json:"maxuploadspeed,string"`
	Online           int    `json:"online,string"`
	UploadTotal      int    `json:"upload,string"`
	UpSpeed          int    `json:"upspeed,string"`
}

type Dev struct {
	Password string `json:"password"`
	Key      string `json:"key"`
	IP       string `json:"ip"`
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
	defer func() {
		sqlDB, err := db.DB()
		checkErr(err)
		sqlDB.Close()
	}()
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
	var (
		cpu       int
		cpu_tp    int
		mem       int
		upSpeed   int
		downSpeed int
		upTotal   int
		downTotal int
		deviceNum int
		devs      []interface{}
		mac       string
	)
	for i, d := range dev {
		ip := d.IP
		routerNum := i
		cpu, cpu_tp, mem, upSpeed, downSpeed, upTotal, downTotal, deviceNum, devs = getRouterStats(i, tokens, ip)
		var count int64
		db.Model(&RouterHistory{}).Where("router_num = ?", routerNum).Count(&count)
		if count >= int64(maxsaved) {
			logrus.Debug("删除历史数据")
			db.Exec("DELETE FROM router_histories WHERE router_num = ? AND created_at = (SELECT MIN(created_at) FROM router_histories WHERE router_num = ? );", routerNum, routerNum)
		}
		db.Create(&RouterHistory{
			Ip:        ip,
			RouterNum: routerNum,
			Cpu:       int(cpu),
			Cpu_tp:    cpu_tp,
			Mem:       int(mem),
			UpSpeed:   int(upSpeed),
			DownSpeed: int(downSpeed),
			UpTotal:   int(upTotal),
			DownTotal: int(downTotal),
			DeviceNum: deviceNum,
		})

		for _, dev := range devs {
			devMap := dev.(map[string]interface{})

			macVal, ok := devMap["mac"]
			if !ok || macVal == "" {
				continue
			}

			data, err := json.Marshal(devMap)
			// logrus.Debug("data: ", string(data))
			checkErr(err)

			var info DeviceInfo
			err = json.Unmarshal(data, &info)
			checkErr(err)
			mac = info.Mac
			upSpeed = int(info.UpSpeed)
			downSpeed = int(info.DownSpeed)
			upTotal = int(info.UploadTotal)
			downTotal = int(info.DownloadTotal)
			db.Create(&DevicesHistory{
				Mac:       mac,
				UpSpeed:   int(upSpeed),
				DownSpeed: int(downSpeed),
				UpTotal:   int(upTotal),
				DownTotal: int(downTotal),
			})
			db.Model(&DevicesHistory{}).Where("mac = ?", mac).Count(&count)
			if count >= int64(maxsaved) {
				logrus.Debug("删除历史数据")
				db.Exec("DELETE FROM devices_histories WHERE mac = ? AND created_at = (SELECT MIN(created_at) FROM devices_histories WHERE mac = ? );", mac, mac)
			}
		}
	}
	defer func() {
		sqlDB, err := db.DB()
		checkErr(err)
		sqlDB.Close()
	}()
}

func GetRouterHistory(databasePath string, routernum int, fixupfloat bool) []RouterHistory {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	checkErr(err)
	var history []RouterHistory
	db.Where("router_num = ?", routernum).Find(&history)
	defer func() {
		sqlDB, err := db.DB()
		checkErr(err)
		sqlDB.Close()
	}()
	return history
}

func GetDeviceHistory(databasePath string, deviceMac string, fixupfloat bool) []DevicesHistory {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	checkErr(err)
	var history []DevicesHistory
	db.Where("mac = ?", deviceMac).Find(&history)
	defer func() {
		sqlDB, err := db.DB()
		checkErr(err)
		sqlDB.Close()
	}()
	return history
}

func getRouterStats(routernum int, tokens map[int]string, ip string) (int, int, int, int, int, int, int, int, []interface{}) {
	if tokens[routernum] == "" {
		return 0, 0, 0, 0, 0, 0, 0, 0, []interface{}{}
	}
	// 需要添加token过期处理
	url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/misystem/status", ip, tokens[routernum])
	resp, err := http.Get(url)
	checkErr(err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	upspeedTemp, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["upspeed"].(string), 64)
	downspeedTemp, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["downspeed"].(string), 64)
	uploadtotalTemp, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["upload"].(string), 64)
	downloadtotalTemp, _ := strconv.ParseFloat(result["wan"].(map[string]interface{})["download"].(string), 64)
	upspeed := int(upspeedTemp)
	downspeed := int(downspeedTemp)
	uploadtotal := int(uploadtotalTemp)
	downloadtotal := int(downloadtotalTemp)

	cpuload := int(result["cpu"].(map[string]interface{})["load"].(float64) * 100)
	cpu_tp := int(result["temperature"].(float64))
	memusage := int(result["mem"].(map[string]interface{})["usage"].(float64) * 100)
	devicenum_now := int(result["count"].(map[string]interface{})["online"].(float64))
	devs := result["dev"].([]interface{})

	return cpuload, cpu_tp, memusage, upspeed, downspeed, uploadtotal, downloadtotal, devicenum_now, devs
}

// func roundToOneDecimal(num int) int {
// 	return math.Round(num*100) / 100
// }

func checkErr(err error) {
	if err != nil {
		logrus.Debug(err)
	}
}
func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0] / 100
}

// func round(val int, roundOn int, places int) (newVal int) {
// 	var rounder int
// 	pow := math.Pow(10, int(places))
// 	intermed := val * pow
// 	_, frac := math.Modf(intermed)

// 	if frac >= roundOn {
// 		rounder = math.Ceil(intermed)
// 	} else {
// 		rounder = math.Floor(intermed)
// 	}
// 	newVal = rounder / pow
// 	return
// }
