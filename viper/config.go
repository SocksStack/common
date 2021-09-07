package viper

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/gobuffalo/packr"
	"os"
)

var Config = viper.New()

func Init(path string) {
	box := packr.NewBox(path)
	configType := "yml"
	defaultConfig, err := box.Find("default.yml")
	if err != nil {
		panic(err)
	}
	Config.SetConfigType(configType)
	err = Config.ReadConfig(bytes.NewReader(defaultConfig))
	if err != nil {
		return
	}

	configs := Config.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		Config.SetDefault(k, v)
	}
	// 在active配置中读取
	active := Config.Get("active")
	if active.(string) != "" {
		activeConfig, err := box.Find(fmt.Sprintf("%s.%s", active, configType))
		if err != nil {
			return
		}
		Config.SetConfigType(configType)
		err = Config.ReadConfig(bytes.NewReader(activeConfig))
	}
	// 在命令行中读取
	env := os.Getenv("MODE")
	if env != "" {
		envConfig, err := box.Find(fmt.Sprintf("%s.%s", env, configType))
		if err != nil {
			return
		}
		Config.SetConfigType(configType)
		err = Config.ReadConfig(bytes.NewReader(envConfig))
	}
}