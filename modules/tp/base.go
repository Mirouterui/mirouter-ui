package tp

import (
	"encoding/json"
	"main/modules/config"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	hardware string
	dev      []config.Dev
	cpu_cmd  *exec.Cmd
	w24g_cmd *exec.Cmd
	w5g_cmd  *exec.Cmd
)

// 获取温度
func GetTemperature(c echo.Context, routernum int, hardware string) (bool, string, string, string, string) {
	if !dev[routernum].RouterUnit {
		return false, "-233", "-233", "-233", "-233"
	}
	var cpu_out, w24g_out, w5g_out []byte
	var err1, err2, err3 error
	var cpu_tp, fanspeed, w24g_tp, w5g_tp string
	switch hardware {
	case "CR8809":
		cpu_cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
		w5g_cmd = exec.Command("cat", "/sys/class/ieee80211/phy0/device/net/wifi0/thermal/temp") //不知道是不是
		cpu_out, err1 = cpu_cmd.Output()
		w5g_out, err3 = w5g_cmd.Output()
		cpu_tp = string(cpu_out)
		fanspeed = "-233"
		w24g_tp = "-233"
		w5g_tp = string(w5g_out)
	case "RA69":
		cpu_cmd = exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
		w24g_cmd = exec.Command("cat", "/sys/class/ieee80211/phy0/device/net/wifi1/thermal/temp")
		w5g_cmd = exec.Command("cat", "/sys/class/ieee80211/phy0/device/net/wifi0/thermal/temp")
		cpu_out, err1 = cpu_cmd.Output()
		w24g_out, err2 = w24g_cmd.Output()
		w5g_out, err3 = w5g_cmd.Output()

		cpu_tp = string(cpu_out)
		fanspeed = "-233"
		w24g_tp = string(w24g_out)
		w5g_tp = string(w5g_out)
	case "R1D":
		type Ubus_data struct {
			Fanspeed    string `json:"fanspeed"`
			Temperature string `json:"temperature"`
		}
		cpu_cmd = exec.Command("ubus", "call", "rmonitor", "status")
		cpu_out, err1 = cpu_cmd.Output()
		var data Ubus_data
		err := json.Unmarshal(cpu_out, &data)
		if err != nil {
			logrus.Error("获取温度失败,报错信息为" + err.Error())
		}
		cpu_tp = data.Temperature
		fanspeed = data.Fanspeed
		w24g_tp = "-233"
		w5g_tp = "-233"
	default:
		return false, "-233", "-233", "-233", "-233"
	}

	if err1 != nil || err2 != nil || err3 != nil {
		logrus.Error("获取温度失败,报错信息为" + err1.Error() + err2.Error() + err3.Error())
	}

	return true, cpu_tp, fanspeed, w24g_tp, w5g_tp
}
