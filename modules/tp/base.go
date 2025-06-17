package tp

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/Mirouterui/mirouter-ui/modules/config"

	"github.com/sirupsen/logrus"
)

var (
	cpuCmd  *exec.Cmd
	w24gCmd *exec.Cmd
	w5gCmd  *exec.Cmd
)

// TemperatureData 存储温度数据
type TemperatureData struct {
	CPU      int `json:"cpu"`
	FanSpeed int `json:"fanspeed"`
	W24G     int `json:"w24g"`
	W5G      int `json:"w5g"`
}

// TemperatureStatus 存储温度数据的状态
type TemperatureStatus struct {
	CPU      bool `json:"cpu"`
	FanSpeed bool `json:"fanspeed"`
	W24G     bool `json:"w24g"`
	W5G      bool `json:"w5g"`
}

// TemperatureResult 返回的温度结果
type TemperatureResult struct {
	Success bool              `json:"success"`
	Data    TemperatureData   `json:"data"`
	Status  TemperatureStatus `json:"status"`
}

// GetTemperature retrieves temperature and fan speed data from a local router device based on its hardware type.
//
// It executes hardware-specific system commands to obtain CPU temperature, fan speed, and WiFi radio temperatures (2.4GHz and 5GHz), returning the results and individual success statuses in a TemperatureResult struct. If the device is not local or the hardware type is unsupported, all values are zero and statuses are false.
func GetTemperature(c interface{}, routerNum int, hardware string, dev []config.Dev) TemperatureResult {
	result := TemperatureResult{
		Success: false,
		Data: TemperatureData{
			CPU:      0,
			FanSpeed: 0,
			W24G:     0,
			W5G:      0,
		},
		Status: TemperatureStatus{
			CPU:      false,
			FanSpeed: false,
			W24G:     false,
			W5G:      false,
		},
	}

	if !dev[routerNum].IsLocal {
		return result
	}

	var cpuOut, w24gOut, w5gOut []byte
	var cpuErr, w24gErr, w5gErr error
	var cpuTemp, fanSpeed, w24gTemp, w5gTemp string

	switch hardware {
	case "CR8809":
		cpuCmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
		w5gCmd = exec.Command("cat", "/sys/class/ieee80211/phy0/device/net/wifi0/thermal/temp")
		cpuOut, cpuErr = cpuCmd.Output()
		w5gOut, w5gErr = w5gCmd.Output()
		cpuTemp = string(cpuOut)
		fanSpeed = "0"
		w24gTemp = "0"
		w5gTemp = string(w5gOut)
	case "RA69":
		cpuCmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
		w24gCmd = exec.Command("cat", "/sys/class/ieee80211/phy0/device/net/wifi1/thermal/temp")
		w5gCmd = exec.Command("cat", "/sys/class/ieee80211/phy0/device/net/wifi0/thermal/temp")
		cpuOut, cpuErr = cpuCmd.Output()
		w24gOut, w24gErr = w24gCmd.Output()
		w5gOut, w5gErr = w5gCmd.Output()

		cpuTpTemp, _ := strconv.Atoi(strings.TrimSpace(string(cpuOut)))
		cpuTemp = strconv.Itoa(cpuTpTemp / 1000)
		fanSpeed = "0"
		w24gTemp = string(w24gOut)
		w5gTemp = string(w5gOut)
	case "R1D":
		type ubusData struct {
			Fanspeed    string `json:"fanspeed"`
			Temperature string `json:"temperature"`
		}
		cpuCmd = exec.Command("ubus", "call", "rmonitor", "status")
		cpuOut, cpuErr = cpuCmd.Output()
		var data ubusData
		err := json.Unmarshal(cpuOut, &data)
		if err != nil {
			logrus.Error("Failed to get temperature, error message: " + err.Error())
		}
		cpuTemp = data.Temperature
		fanSpeed = data.Fanspeed
		w24gTemp = "0"
		w5gTemp = "0"
	case "RB06":
		cpuCmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
		w24gCmd = exec.Command("iwpriv", "wl0", "stat", "|", "grep", "CurrentTemperature")
		w5gCmd = exec.Command("iwpriv", "wl1", "stat", "|", "grep", "CurrentTemperature")
		cpuOut, cpuErr = cpuCmd.Output()
		w24gOut, w24gErr = w24gCmd.Output()
		w5gOut, w5gErr = w5gCmd.Output()

		cpuTemp = string(cpuOut)
		re := regexp.MustCompile(`CurrentTemperature\s*=\s*(\d+)`)
		matches24g := re.FindStringSubmatch(string(w24gOut))
		if len(matches24g) > 1 {
			w24gTemp = matches24g[1]
		} else {
			w24gTemp = "0"
		}
		matches5g := re.FindStringSubmatch(string(w5gOut))
		if len(matches5g) > 1 {
			w5gTemp = matches5g[1]
		} else {
			w5gTemp = "0"
		}
		fanSpeed = "0"
	default:
		return result
	}

	// 处理错误并记录日志
	if cpuErr != nil || w24gErr != nil || w5gErr != nil {
		logrus.Error("Failed to get temperature, error message: " + cpuErr.Error() + w24gErr.Error() + w5gErr.Error())
	}

	// 转换温度数据并更新状态
	if cpuTemp != "" {
		result.Data.CPU, _ = strconv.Atoi(strings.ReplaceAll(cpuTemp, "\n", ""))
		result.Status.CPU = cpuErr == nil
	}

	if fanSpeed != "" {
		result.Data.FanSpeed, _ = strconv.Atoi(strings.ReplaceAll(fanSpeed, "\n", ""))
		result.Status.FanSpeed = true
	}

	if w24gTemp != "" {
		result.Data.W24G, _ = strconv.Atoi(strings.ReplaceAll(w24gTemp, "\n", ""))
		result.Status.W24G = w24gErr == nil
	}

	if w5gTemp != "" {
		result.Data.W5G, _ = strconv.Atoi(strings.ReplaceAll(w5gTemp, "\n", ""))
		result.Status.W5G = w5gErr == nil
	}

	result.Success = true
	return result
}
