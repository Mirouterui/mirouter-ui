package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Mirouterui/mirouter-ui/modules/download"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Dev struct {
	Password   string `mapstructure:"password"`
	Key        string `mapstructure:"key"`
	IP         string `mapstructure:"ip"`
	RouterUnit bool   `mapstructure:"routerunit"`
	IsLocal    bool   `mapstructure:"islocal"`
}

type History struct {
	Enable     bool `mapstructure:"enable"`
	MaxDeleted int  `mapstructure:"maxsaved"`
	Sampletime int  `mapstructure:"sampletime"`
}

type AppConfig struct {
	Dev               []Dev   `mapstructure:"dev"`
	History           History `mapstructure:"history"`
	Debug             bool    `mapstructure:"debug"`
	Port              int     `mapstructure:"port"`
	Address           string  `mapstructure:"address"`
	Tiny              bool    `mapstructure:"tiny"`
	FlushTokenTime    int     `mapstructure:"flushTokenTime"`
	Netdata_routernum int     `mapstructure:"netdata_routernum"`
	Workdirectory     string  `mapstructure:"-"`
	Databasepath      string  `mapstructure:"-"`
	ApiKey            string  `mapstructure:"api_key"`
	SafeMode          bool    `mapstructure:"safemode"`
}

var (
	// Global configuration instance
	Cfg *AppConfig

	// Command line parameters
	configPath      string
	workdirectory   string
	databasepath    string
	autocheckupdate string
)

func init() {
	appPath, _ := os.Executable()
	flag.StringVar(&configPath, "config", filepath.Join(filepath.Dir(appPath), "config.yaml"), "configuration file path")
	flag.StringVar(&workdirectory, "workdirectory", "", "working directory path")
	flag.StringVar(&databasepath, "databasepath", filepath.Join(filepath.Dir(appPath), "database.db"), "database path")
	flag.StringVar(&autocheckupdate, "autocheckupdate", "true", "auto check updates")
	flag.Parse()
}

func LoadConfig() (*AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set default values
	v.SetDefault("dev", []Dev{
		{
			Password:   "",
			Key:        "a2ffa5c9be07488bbb04a3a47d3c5f6a",
			IP:         "192.168.31.1",
			RouterUnit: false,
			IsLocal:    false,
		},
	})
	v.SetDefault("history.enable", false)
	v.SetDefault("history.maxsaved", 3000)
	v.SetDefault("history.sampletime", 86400)
	v.SetDefault("debug", true)
	v.SetDefault("port", 6789)
	v.SetDefault("tiny", false)
	v.SetDefault("flushTokenTime", 1800)
	v.SetDefault("netdata_routernum", 0)
	v.SetDefault("api_key", "")
	v.SetDefault("safemode", true)
	v.SetDefault("address", "0.0.0.0")

	// Support environment variable overrides
	v.AutomaticEnv()
	v.SetEnvPrefix("MIROUTERUI")

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := v.WriteConfigAs(configPath); err != nil {
				return nil, fmt.Errorf("failed to create default config file: %w", err)
			}
			logrus.Info("default config file created, please modify it and restart")
			time.Sleep(5 * time.Second)
			os.Exit(0)
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	Cfg = &AppConfig{}
	if err := v.Unmarshal(Cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	Cfg.Workdirectory = workdirectory
	Cfg.Databasepath = databasepath

	if len(Cfg.Dev) == 0 {
		return nil, fmt.Errorf("router information not filled, please check config file")
	}

	if Cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Check for updates
	autocheckupdatebool, _ := strconv.ParseBool(autocheckupdate)
	if !Cfg.Tiny {
		download.DownloadStatic(Cfg.Workdirectory, false, autocheckupdatebool)
	}

	return Cfg, nil
}
