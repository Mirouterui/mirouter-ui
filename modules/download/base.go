package download

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Version string
)

func DownloadStatic(basedirectory string, force bool) error {
	directory := "static"
	if basedirectory != "" {
		directory = filepath.Join(basedirectory, "static")
	}
	if force {
		//删除
		os.RemoveAll(directory)
	}

	_, err := os.Stat(directory)
	if os.IsNotExist(err) || force {
		logrus.Info("正从'Mirouterui/static'下载静态资源")
		downloadfile(directory)
		return nil
	}

	// 读取/static/version/index.html
	f, err := os.Open(filepath.Join(directory, "version", "index.html"))
	if err != nil {
		logrus.Info("无法读取静态资源版本号，重新下载")
		os.RemoveAll(directory)
		downloadfile(directory)
		return err
	}

	defer f.Close()
	forntendVersion, err := io.ReadAll(f)
	checkErr(err)
	logrus.Info("静态资源已存在，版本号为" + string(forntendVersion))

	resp, err := http.Get("https://mrui-api.hzchu.top/checkupdate")

	if err != nil {
		logrus.Info("无法获取更新信息，跳过检查")
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	checkErr(err)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["backversion"] != string(Version) {
		message := fmt.Sprintf("后端程序发现新版本(%v)，请及时更新", result["backversion"])
		logrus.Info(message)
	}

	if result["frontversion"] != string(forntendVersion) {
		message := fmt.Sprintf("前端文件发现新版本(%v)，正在重新下载", result["frontversion"])
		logrus.Info(message)
		os.RemoveAll(directory)
		downloadfile(directory)
	}
	return nil
}
func downloadfile(directory string) {
	resp, err := http.Get("http://mrui-api.hzchu.top/downloadstatic")
	checkErr(err)
	defer resp.Body.Close()

	out, err := os.CreateTemp("", "*.zip")
	checkErr(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	checkErr(err)

	err = unzip(out.Name(), directory)
	checkErr(err)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		logrus.Info("静态资源下载失败，请尝试手动下载")
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		checkErr(err)
		fname := f.Name
		if len(fname) > 26 {
			fname = fname[26:]
		}
		fpath := filepath.Join(dest, fname)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
		}

		rc.Close()
	}

	return nil
}

func checkErr(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}
