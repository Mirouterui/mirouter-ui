package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	password   string
	key        string
	iv         string
	ip         string
	token      string
	debug      bool
	port       int
	routername string
	hardware   string
)

type Config struct {
	Password string `json:"password"`
	Key      string `json:"key"`
	Iv       string `json:"iv"`
	Ip       string `json:"ip"`
	Debug    bool   `json:"debug"`
	Port     int    `json:"port"`
}

func init() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(filepath.Dir(exePath), "config.json")
	fmt.Println(configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	password = config.Password
	key = config.Key
	iv = config.Iv
	ip = config.Ip
	debug = config.Debug
	port = config.Port
	// fmt.Println(password)
	// fmt.Println(key)
	// fmt.Println(iv)

}

func createNonce() string {
	typeVar := 0
	deviceID := "" //无效参数
	timeVar := int(time.Now().Unix())
	randomVar := rand.Intn(10000)
	return fmt.Sprintf("%d_%s_%d_%d", typeVar, deviceID, timeVar, randomVar)
}

func hashPassword(pwd string, nonce string, key string) string {
	pwdKey := pwd + key
	pwdKeyHash := sha1.New()
	pwdKeyHash.Write([]byte(pwdKey))
	pwdKeyHashStr := fmt.Sprintf("%x", pwdKeyHash.Sum(nil))

	noncePwdKey := nonce + pwdKeyHashStr
	noncePwdKeyHash := sha1.New()
	noncePwdKeyHash.Write([]byte(noncePwdKey))
	noncePwdKeyHashStr := fmt.Sprintf("%x", noncePwdKeyHash.Sum(nil))

	return noncePwdKeyHashStr
}
func newhashPassword(pwd string, nonce string, key string) string {
	pwdKey := pwd + key
	pwdKeyHash := sha256.Sum256([]byte(pwdKey))
	pwdKeyHashStr := hex.EncodeToString(pwdKeyHash[:])

	noncePwdKey := nonce + pwdKeyHashStr
	noncePwdKeyHash := sha256.Sum256([]byte(noncePwdKey))
	noncePwdKeyHashStr := hex.EncodeToString(noncePwdKeyHash[:])

	return noncePwdKeyHashStr
}
func getrouterinfo() int {

	// 发送 GET 请求
	ourl := fmt.Sprintf("http://%s/cgi-bin/luci/api/xqsystem/init_info", ip)
	response, err := http.Get(ourl)
	if err != nil {
		return 0
	}
	defer response.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0
	}

	// 解析 JSON
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0
	}
	//提取routername
	routername = data["routername"].(string)
	hardware = data["hardware"].(string)
	fmt.Println("路由器型号为：", hardware)
	fmt.Println("路由器名称为：", routername)
	// 检查 newEncryptMode
	newEncryptMode, ok := data["newEncryptMode"].(float64)
	if !ok {
		debugPrint("使用旧加密模式")
		return 0
	}

	if newEncryptMode != 0 {
		debugPrint("使用新加密模式")
		fmt.Println("当前路由器可能无法正常获取某些数据！")
		return 1
	}
	return 0
}
func updateToken() {
	debugPrint("获取路由器信息...")
	newEncryptMode := getrouterinfo()
	// fmt.Println(newEncryptMode)
	fmt.Println("更新token...")
	nonce := createNonce()
	var hashedPassword string

	if newEncryptMode == 1 {
		hashedPassword = newhashPassword(password, nonce, key)
	} else {
		hashedPassword = hashPassword(password, nonce, key)
	}

	ourl := fmt.Sprintf("http://%s/cgi-bin/luci/api/xqsystem/login", ip)
	params := url.Values{}
	params.Set("username", "admin")
	params.Set("password", hashedPassword)
	params.Set("logtype", "2")
	params.Set("nonce", nonce)

	resp, err := http.PostForm(ourl, params)
	if err != nil {
		fmt.Println("登录失败，请检查配置或路由器状态")
		fmt.Println("30秒后退出程序")
		time.Sleep(30 * time.Second)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	code := int(result["code"].(float64))
	if code == 0 {
		debugPrint("当前token为:" + fmt.Sprint(result["token"]))
		token = result["token"].(string)
	} else {
		fmt.Println("登录失败，请检查配置")
		fmt.Println("30秒后退出程序")
		fmt.Println(string(body))
		time.Sleep(30 * time.Second)
		os.Exit(1)
	}

}
func debugPrint(msg string) {
	if debug == true { // 假设debug是一个全局变量
		fmt.Println(msg) // 如果debug为true，就打印消息
	}
}

func main() {
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/:apipath", func(c echo.Context) error {
		apipath := c.Param("apipath")
		switch apipath {
		case "xqsystem/router_name":
			return c.JSON(http.StatusOK, map[string]interface{}{"code": 0, "routerName": routername})
		case "misystem/status", "misystem/devicelist", "xqsystem/internet_connect", "xqsystem/fac_info", "misystem/messages":
			url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/%s", ip, token, apipath)
			resp, err := http.Get(url)
			if err != nil {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"code": 1101,
					"msg":  "MiRouterのapi调用出错，请检查配置或路由器状态",
				})
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)
			return c.JSON(http.StatusOK, result)
		default:
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1102,
				"msg":  "该api不支持免密调用",
			})
		}
	})

	// var contentHandler = echo.WrapHandler(http.FileServer(http.FS(static)))
	// var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/static/$1"})

	// e.GET("/*", contentHandler, contentRewrite)
	e.Static("/", "static")

	updateToken()
	go func() {
		for range time.Tick(30 * time.Minute) {
			updateToken()
		}
	}()

	e.Start(":" + fmt.Sprint(port))
}
