package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path"
)

var Configs *viper.Viper

// 加载默认配置
func loadDefaultConfig() {
	defaultConfig := map[string]interface{}{
		"# 下面是默认配置文件": nil,
		"logs": map[string]interface{}{
			"# logs.path 日志文件路径":  nil,
			"# logs.level 日志输出级别": nil,
			"path":                "./logs/logrus.log",
			"level":               "info",
		},
		"# goLeaf 配置信息": nil,
		"goLeaf": map[string]interface{}{
			"name": "ip:port",
			"port": "9090",
			"snowflake": map[string]interface{}{
				"enable":           true,
				"start":            "2020-01-01",
				"zookeeperAddress": "",
			},
		},
	}
	Configs = viper.New()
	//将默认值设置到config中
	Configs.MergeConfigMap(defaultConfig)
}

// LoadConfig 加载配置文件
func LoadConfig() {
	loadDefaultConfig()
	configDir := "./config/"
	configName := "config"
	configType := "yaml"
	// 如果 Configs 目录不存在，则创建它
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println(err)
		}
	}
	// 解析配置文件
	Configs.AddConfigPath(configDir)
	Configs.SetConfigName(configName)
	Configs.SetConfigType(configType)

	if err := Configs.ReadInConfig(); err != nil {
		// 如果找不到配置文件，则提醒生成配置文件并创建它
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			file := fmt.Sprintf("%s.%s", configName, configType)
			configPath := path.Join(configDir + file)
			fmt.Printf("[warning] Config file not found. Generating default config file at %s\n", configPath)
			if err := Configs.WriteConfigAs(configPath); err != nil {
				fmt.Printf("[error] Failed to generate default config file. %s", err)
			}
			// 再次读取配置文件
			if err := Configs.ReadInConfig(); err != nil {
				fmt.Printf("[error] Failed to read config file. %s", err)
			}
		} else {
			// 配置文件被找到，但产生了另外的错误
			fmt.Printf("[error] Failed to parse config file. %s", err)
		}
	}
	loadSetting()

	//实时监控配置文件变化
	Configs.WatchConfig()
	Configs.OnConfigChange(func(e fsnotify.Event) {
		reloadConfig()
	})
}
func reloadConfig() {
	// 配置文件发生变更之后会调用的回调函数
	loadSetting()
}
