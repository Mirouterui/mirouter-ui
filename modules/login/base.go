package login

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
	"time"

	"github.com/sirupsen/logrus"
)

var (
	password   string
	key        string
	ip         string
	token      string
	routername string
	hardware   string
)

// createNonce generates a nonce string using a fixed type, a hardcoded MAC address, the current Unix timestamp, and a random integer.
func createNonce() string {
	typeVar := 0
	deviceID := "00:e0:4f:27:3d:09" //MAC
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
// getrouterinfo retrieves router information and determines the encryption mode.
//
// Sends an HTTP request to the router's init_info API endpoint, parses the response for router name and hardware model, and checks if the new encryption mode is enabled.
//
// Returns true if the router uses the new encryption mode, along with the router name and hardware model. Returns false and empty strings if the request or parsing fails.
func getrouterinfo(ip string) (bool, string, string) {

	// 发送 GET 请求
	ourl := fmt.Sprintf("http://%s/cgi-bin/luci/api/xqsystem/init_info", ip)
	response, err := http.Get(ourl)
	if err != nil {
		return false, "", ""
	}
	defer response.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, "", ""
	}

	// 解析 JSON
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, "", ""
	}
	//提取routername
	routername = data["routername"].(string)
	hardware = data["hardware"].(string)
	logrus.Debug("Router model: " + hardware)
	logrus.Debug("Router name: " + routername)
	// 检查 newEncryptMode
	newEncryptMode, ok := data["newEncryptMode"].(float64)
	if !ok {
		logrus.Debug("Using old encryption mode")
		return false, routername, hardware
	}

	if newEncryptMode != 0 {
		logrus.Debug("Using new encryption mode")
		logrus.Info("The current router may not be able to fetch certain data properly!")
		return true, routername, hardware
	}
	return false, routername, hardware
}

// CheckRouterAvailability returns true if the router at the specified IP address is reachable via HTTP within 5 seconds.
func CheckRouterAvailability(ip string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	_, err := client.Get("http://" + ip)
	if err != nil {
		logrus.Info("Router " + ip + " is not available, please check configuration or router status")
		return false
	}

	return true
}
// GetToken authenticates with the router at the specified IP address and retrieves an authentication token, router name, and hardware information.
// If the router is unavailable or authentication fails, the function logs the error, waits 5 seconds, and exits the program.
// Returns the authentication token, router name, and hardware string on success. If the router is unavailable, returns an empty token and an error message.
func GetToken(password string, key string, ip string) (string, string, string) {
	logrus.Debug("Checking router availability...")
	if !CheckRouterAvailability(ip) {
		return "", "Router is not available", ""
	}
	logrus.Debug("Getting router information...")
	newEncryptMode, routername, hardware := getrouterinfo(ip)
	logrus.Info("Updating token...")
	nonce := createNonce()
	var hashedPassword string

	if newEncryptMode {
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
		logrus.Info("Login failed, please check configuration or router status")
		logrus.Info(err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	var code int
	if result["code"] != nil {
		code = int(result["code"].(float64))
	} else {
		logrus.Info("Router login request returned empty! Please check configuration")
	}

	if code == 0 {
		logrus.Debug("Current token: " + fmt.Sprint(result["token"]))
		token = result["token"].(string)
	} else {
		logrus.Info("Login failed, please check configuration, the following is the return output:")
		logrus.Info(string(body))
		logrus.Info("Exiting program in 5 seconds")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	return token, routername, hardware
}
