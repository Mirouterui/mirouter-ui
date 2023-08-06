package main

import (
	"crypto/sha1"
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	password string
	key      string
	iv       string
	ip       string
	token    string
)

//go:embed static/*
var static embed.FS

func init() {
	data, err := ioutil.ReadFile("config.txt")
	if err != nil {
		panic(err)
	}
	config := string(data)
	password = strings.Split(strings.Split(config, `password = "`)[1], `"`)[0]
	key = strings.Split(strings.Split(config, `key = "`)[1], `"`)[0]
	iv = strings.Split(strings.Split(config, `iv = "`)[1], `"`)[0]
	ip = strings.Split(strings.Split(config, `ip = "`)[1], `"`)[0]
	// fmt.Println(password)
	// fmt.Println(key)
	// fmt.Println(iv)

}

func createNonce() string {
	typeVar := 0
	deviceID := "2a:45:6f:11:0b:a5"
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
func updateToken() {
	fmt.Println("更新token...")
	nonce := createNonce()
	hashedPassword := hashPassword(password, nonce, key)
	ourl := fmt.Sprintf("http://%s/cgi-bin/luci/api/xqsystem/login", ip)
	params := url.Values{}
	params.Set("username", "admin")
	params.Set("password", hashedPassword)
	params.Set("logtype", "2")
	params.Set("nonce", nonce)

	resp, err := http.PostForm(ourl, params)
	if err != nil {
		fmt.Println("登录失败，请检查配置或路由器状态")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	// if result["code"] != 0 {
	// 	fmt.Println("登录失败，请检查配置")
	// 	return
	// }
	json.Unmarshal(body, &result)
	// fmt.Println(result)
	token = result["token"].(string)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/:apipath", func(c echo.Context) error {
		url := fmt.Sprintf("http://%s/cgi-bin/luci/;stok=%s/api/%s", ip, token, c.Param("apipath"))
		resp, err := http.Get(url)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 1101,
				"msg":  "MiRouterのapi调用出错，请检查配置或路由器状态",
			})
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		return c.JSON(http.StatusOK, result)
	})

	var contentHandler = echo.WrapHandler(http.FileServer(http.FS(static)))
	var contentRewrite = middleware.Rewrite(map[string]string{"/*": "/static/$1"})
	
	e.GET("/*", contentHandler, contentRewrite)


	updateToken()
	go func() {
		for range time.Tick(30 * time.Minute) {
			updateToken()
		}
	}()

	e.Start(":6789")
}
