package netdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Data struct {
	API              int      `json:"api"`
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	ViewUpdateEvery  int      `json:"view_update_every"`
	UpdateEvery      int      `json:"update_every"`
	FirstEntry       int      `json:"first_entry"`
	LastEntry        int      `json:"last_entry"`
	Before           int      `json:"before"`
	After            int      `json:"after"`
	DimensionNames   []string `json:"dimension_names"`
	DimensionIDs     []string `json:"dimension_ids"`
	LatestValues     []int    `json:"latest_values"`
	ViewLatestValues []int    `json:"view_latest_values"`
	Dimensions       int      `json:"dimensions"`
	Points           int      `json:"points"`
	Format           string   `json:"format"`
	Result           Result   `json:"result"`
	Min              int      `json:"min"`
	Max              int      `json:"max"`
}

type CacheData struct {
	CpuLoad       int
	MemAvailable  int
	MemTotal      int
	MemUsage      int
	UpSpeed       int
	DownSpeed     int
	Temperature   int
	Deviceonline  int
	uploadtotal   int
	downloadtotal int
}
type Result struct {
	Labels []string `json:"labels"`
	Data   [][]int  `json:"data"`
}
type Dimension struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type DataForAllMetrics struct {
	Name        string               `json:"name"`
	Family      string               `json:"family"`
	Context     string               `json:"context"`
	Units       string               `json:"units"`
	LastUpdated int                  `json:"last_updated"`
	Dimensions  map[string]Dimension `json:"dimensions"`
}

var c = cache.New(2*time.Second, 4*time.Second)

func ProcessData(ip string, token string) (int, int, int, int, int, int, int, int, int, int) {
	cacheKey := "netdata-cache"
	if x, found := c.Get(cacheKey); found {
		data := x.(CacheData)
		return data.CpuLoad, data.MemAvailable, data.MemTotal, data.MemUsage, data.UpSpeed, data.DownSpeed, data.Temperature, data.Deviceonline, data.uploadtotal, data.downloadtotal
	}
	url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/misystem/status", ip, token)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Info(err)
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	cpuLoad, _ := strconv.Atoi(fmt.Sprintf("%.0f", result["cpu"].(map[string]interface{})["load"].(float64)*100))
	memUsage, _ := strconv.Atoi(fmt.Sprintf("%.0f", result["mem"].(map[string]interface{})["usage"].(float64)*100))
	memTotal, _ := strconv.Atoi(result["mem"].(map[string]interface{})["total"].(string)[:len(result["mem"].(map[string]interface{})["total"].(string))-2])
	upSpeed, _ := strconv.Atoi(result["wan"].(map[string]interface{})["upspeed"].(string))
	downSpeed, _ := strconv.Atoi(result["wan"].(map[string]interface{})["downspeed"].(string))
	uploadtotal, _ := strconv.Atoi(result["wan"].(map[string]interface{})["upload"].(string))
	downloadtotal, _ := strconv.Atoi(result["wan"].(map[string]interface{})["download"].(string))
	upSpeed = upSpeed / 1024 * 8
	downSpeed = downSpeed / 1024 * 8

	temperature, _ := strconv.Atoi(fmt.Sprintf("%.0f", result["temperature"].(float64)))
	memAvailable := memTotal * (100 - memUsage) / 100
	deviceonline := int(result["count"].(map[string]interface{})["online"].(float64))

	data := CacheData{
		CpuLoad:       cpuLoad,
		MemAvailable:  memAvailable,
		MemTotal:      memTotal,
		MemUsage:      memUsage,
		UpSpeed:       upSpeed,
		DownSpeed:     downSpeed,
		Temperature:   temperature,
		Deviceonline:  deviceonline,
		uploadtotal:   uploadtotal,
		downloadtotal: downloadtotal,
	}
	c.Set(cacheKey, data, cache.DefaultExpiration)
	return cpuLoad, memAvailable, memTotal, memUsage, upSpeed, downSpeed, temperature, deviceonline, uploadtotal, downloadtotal
}

func GenerateArray(id string, latestValue int, FirstEntry int, dimensionName string, dimensionID string) Data {
	time := int(time.Now().Unix())
	return Data{
		API:              1,
		ID:               id,
		Name:             id,
		ViewUpdateEvery:  2,
		UpdateEvery:      1,
		FirstEntry:       FirstEntry,
		LastEntry:        time,
		Before:           time - 1,
		After:            time - 2,
		DimensionNames:   []string{dimensionName},
		DimensionIDs:     []string{dimensionID},
		LatestValues:     []int{latestValue},
		ViewLatestValues: []int{latestValue},
		Dimensions:       1,
		Points:           1,
		Format:           "json",
		Result: Result{
			Labels: []string{"time", dimensionName},
			Data:   [][]int{{1703897272, latestValue}},
		},
		Min: latestValue,
		Max: latestValue,
	}
}

func GenerateDataForAllMetrics(id string, family string, units string, latestValue int, dimensionName string) DataForAllMetrics {
	time := int(time.Now().Unix())
	return DataForAllMetrics{
		Name:        id,
		Family:      family,
		Context:     "",
		Units:       units,
		LastUpdated: time,
		Dimensions: map[string]Dimension{
			dimensionName: {
				Name:  dimensionName,
				Value: latestValue,
			},
		},
	}
}
