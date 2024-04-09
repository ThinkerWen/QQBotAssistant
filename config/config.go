package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"os"
)

var CONFIG = initConfig()

var ApiUrl = CONFIG.ApiUrl
var Molly = CONFIG.Molly
var HeroPower = CONFIG.HeroPower
var Sensitive = CONFIG.Sensitive

type Config struct {
	Name      string          `mapstructure:"name"`
	ApiUrl    string          `mapstructure:"api_url"`
	Hosts     []int64         `mapstructure:"hosts"`
	HeroPower HeroPowerConfig `mapstructure:"hero_power" json:"hero_power"`
	Molly     MollyConfig     `mapstructure:"molly"`
	Sensitive SensitiveConfig `mapstructure:"sensitive"`
}

type HeroPowerConfig struct {
	Enable bool    `mapstructure:"enable" json:"enable"`
	Token  string  `mapstructure:"token" json:"token"`
	Groups []int64 `mapstructure:"groups"`
}

type MollyConfig struct {
	Enable    bool    `mapstructure:"enable"`
	QQ        int64   `mapstructure:"qq"`
	Name      string  `mapstructure:"name"`
	ApiKey    string  `mapstructure:"api_key" json:"api_key"`
	ApiSecret string  `mapstructure:"api_secret" json:"api_secret"`
	Groups    []int64 `mapstructure:"groups"`
}

type SensitiveConfig struct {
	Enable      bool    `mapstructure:"enable"`
	Token       string  `mapstructure:"token"`
	AlertTimes  int     `mapstructure:"alert_times" json:"alert_times"`
	ShutSeconds int     `mapstructure:"shut_seconds" json:"shut_seconds"`
	Groups      []int64 `mapstructure:"groups"`
}

func initConfig() Config {
	workDir, _ := os.Getwd()
	viper.AddConfigPath(workDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("application")
	initDefaultConfig()
	_ = viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("读取配置文件失败", err)
	}
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Println("解析结构体失败", err)
	}
	return config
}

func initDefaultConfig() {
	var data []byte
	var config map[string]interface{}

	heroPower := new(HeroPowerConfig)
	heroPower.Token = "free"
	heroPower.Enable = true
	viper.SetDefault("hero_power", *heroPower)

	molly := new(MollyConfig)
	molly.Enable = true
	molly.Groups = make([]int64, 0)
	data, _ = json.Marshal(*molly)
	_ = json.Unmarshal(data, &config)
	viper.SetDefault("molly", config)

	data = []byte{}
	config = make(map[string]interface{})
	sensitive := new(SensitiveConfig)
	sensitive.Token = "free"
	sensitive.Enable = true
	sensitive.AlertTimes = 3
	sensitive.ShutSeconds = 60
	sensitive.Groups = make([]int64, 0)
	data, _ = json.Marshal(*sensitive)
	_ = json.Unmarshal(data, &config)
	viper.SetDefault("sensitive", config)

	viper.SetDefault("hosts", []int64{})
	viper.SetDefault("name", "QQBotAssistant")
	viper.SetDefault("api_url", "http://127.0.0.1:8086")
}

func ReLoadSubConfig() {
	Molly = CONFIG.Molly
	HeroPower = CONFIG.HeroPower
	Sensitive = CONFIG.Sensitive
}
