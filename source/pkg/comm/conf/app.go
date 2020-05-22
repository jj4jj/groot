package conf

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var appConfig AppConfig

func LoadConfig(cfgPath string) error {

	fmt.Println("Load config file path:", cfgPath)
	//init config
	if e := initConfig(cfgPath); e != nil {
		return e
	}

	// 初始化日志包
	log.InitWithConfig(&appConfig.LogConf)

	// 监控配置文件变化并热加载程序
	watchConfig()

	return nil
}

func SaveConfig(cfgPath string) error {
	//init config
	var bjson []byte
	var er error
	fmt.Println("save config json text path: ", cfgPath)
	if bjson, er = json.MarshalIndent(*DefaultAppConfig(), "", "    "); er != nil {
		fmt.Fprintf(os.Stderr, "save config path error with file path:%s\n", cfgPath)
		return er
	}

	ioutil.WriteFile(cfgPath, bjson, 0777)

	return nil
}

func GetBool(name string) bool {
	return viper.GetBool(name)
}
func GetInt(name string) int {
	return viper.GetInt(name)
}
func GetString(name string) string {
	return viper.GetString(name)
}

func GetAppConfig() *AppConfig {
	return &appConfig
}

func initConfig(configPath string) error {
	if configPath != "" {
		viper.SetConfigFile(configPath) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml") //设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	viper.SetEnvPrefix("DEMO")  // 读取环境变量的前缀为DEMO
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if e := viper.ReadInConfig(); e != nil {
		return e
	}
	if e := viper.Unmarshal(&appConfig); e != nil {
		return e
	}

	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %s ", e.Name)
	})
}
