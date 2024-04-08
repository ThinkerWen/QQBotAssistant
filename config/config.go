package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var CONFIG = initConfig()

var ApiUrl = CONFIG.ApiUrl
var HeroPower = CONFIG.HeroPower
var Molly = CONFIG.Molly
var Sensitive = CONFIG.Sensitive

type Config struct {
	Name      string          `mapstructure:"name"`
	ApiUrl    string          `mapstructure:"api_url"`
	HeroPower HeroPowerConfig `mapstructure:"hero_power"`
	Molly     MollyConfig     `mapstructure:"molly"`
	Sensitive SensitiveConfig `mapstructure:"sensitive"`
}

type HeroPowerConfig struct {
	Enable bool    `mapstructure:"enable"`
	Token  string  `mapstructure:"token"`
	Hosts  []int64 `mapstructure:"hosts"`
	Groups []int64 `mapstructure:"groups"`
}

type MollyConfig struct {
	Enable    bool   `mapstructure:"enable"`
	QQ        int64  `mapstructure:"qq"`
	Name      string `mapstructure:"name"`
	ApiKey    string `mapstructure:"api_key"`
	ApiSecret string `mapstructure:"api_secret"`
}

type SensitiveConfig struct {
	Enable bool   `mapstructure:"enable"`
	Token  string `mapstructure:"token"`
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
	heroPower := new(HeroPowerConfig)
	heroPower.Token = "free"
	heroPower.Enable = true
	viper.SetDefault("hero_power", *heroPower)

	molly := new(MollyConfig)
	molly.Enable = true
	viper.SetDefault("molly", *molly)

	sensitive := new(SensitiveConfig)
	sensitive.Token = "free"
	sensitive.Enable = true
	viper.SetDefault("molly", *molly)

	viper.SetDefault("api_url", "http://127.0.0.1:8086")
}

func ReLoadSubConfig() {
	HeroPower = CONFIG.HeroPower
	Molly = CONFIG.Molly
	Sensitive = CONFIG.Sensitive
}
