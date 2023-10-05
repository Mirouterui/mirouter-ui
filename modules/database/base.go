package database

import (
	"io"
	"math"
	"strconv"

	"encoding/json"
	"fmt"
	"main/modules/config"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// For database
type History struct {
	gorm.Model
	Ip        string
	RouterNum int
	CPU       float64
	Cpu_tp    int
	Mem       float64
	UpSpeed   float64
	DownSpeed float64
	UpTotal   float64
	DownTotal float64
	DeviceNum int
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
// - databasepath: a string variable that holds the path to the SQLite database file.
//
// Returns: None.
func CheckDatabase(databasepath string) {
	// databasepath is a variable that holds the path to the SQLite database file
	db, err := gorm.Open(sqlite.Open(databasepath), &gorm.Config{})
	checkErr(err)

	// Check if the history table exists, if not, create it
	err = db.AutoMigrate(&History{})
	checkErr(err)

	// Perform CRUD operations on the history table using db.Create, db.First, db.Update, db.Delete methods
}

// Savetodb saves device statistics to the database.
//
// Parameters:
// - databasepath: the path to the database.
// - dev: an array of device configurations.
// - tokens: a map of token IDs to strings.
// - maxdeleted: the maximum number of records to delete.
func Savetodb(databasepath string, dev []config.Dev, tokens map[int]string, maxdeleted int64) {
	db, err := gorm.Open(sqlite.Open(databasepath), &gorm.Config{})
	checkErr(err)
	for i, d := range dev {
		ip := d.IP
		routerNum := i
		cpu, cpu_tp, mem, upSpeed, downSpeed, upTotal, downTotal, deviceNum := getDeviceStats(i, tokens, ip)
		var count int64
		db.Model(&History{}).Where("router_num = ?", routerNum).Count(&count)
		if count >= maxdeleted {
			logrus.Debug("删除历史数据")
			db.Exec("DELETE FROM histories WHERE router_num = ? AND created_at = (SELECT MIN(created_at) FROM histories WHERE router_num = ? );", routerNum, routerNum)

		}
		db.Create(&History{
			Ip:        ip,
			RouterNum: routerNum,
			CPU:       cpu,
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

func Getdata(databasepath string, routernum int) []History {
	db, err := gorm.Open(sqlite.Open(databasepath), &gorm.Config{})
	checkErr(err)
	var history []History
	db.Where("router_num = ?", routernum).Find(&history)
	return history
}

// getDeviceStats retrieves the device statistics from the specified router.
//
// Parameters:
// - routernum: The router number.
// - tokens: A map containing the tokens.
// - ip: The IP address of the router.
//
// Returns:
// - cpuload: The CPU load.
// - cpu_tp: The CPU temperature.
// - memusage: The memory usage.
// - upspeed: The upload speed.
// - downspeed: The download speed.
// - uploadtotal: The total upload amount.
// - downloadtotal: The total download amount.
// - devicenum_now: The number of online devices.
func getDeviceStats(routernum int, tokens map[int]string, ip string) (float64, int, float64, float64, float64, float64, float64, int) {
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
	cpuload := roundToOneDecimal(result["cpu"].(map[string]interface{})["load"].(float64) * 100)
	cpu_tp := int(result["temperature"].(float64))
	memusage := roundToOneDecimal(result["mem"].(map[string]interface{})["usage"].(float64) * 100)
	devicenum_now := int(result["count"].(map[string]interface{})["online"].(float64))

	return cpuload, cpu_tp, memusage, upspeed, downspeed, uploadtotal, downloadtotal, devicenum_now
}

func roundToOneDecimal(num float64) float64 {
	return math.Round(num*100) / 100
}

func checkErr(err error) {
	if err != nil {
		logrus.Debug(err)
	}
}
