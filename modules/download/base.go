package download

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
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

	return nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
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
		logrus.Debug(err)
	}
}
