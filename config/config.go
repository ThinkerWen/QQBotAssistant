package config

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
)

var CONFIG = initConfig()

var ApiUrl = CONFIG.ApiUrl
var Molly = CONFIG.Molly
var HeroPower = CONFIG.HeroPower
var Sensitive = CONFIG.Sensitive
var AutoReply = CONFIG.AutoReply
var OnlineCourse = CONFIG.OnlineCourse
var AutoReplyList = make([]map[string]interface{}, 0)

type Config struct {
	Name         string             `mapstructure:"name"`
	ApiUrl       string             `mapstructure:"api_url"`
	Hosts        []int64            `mapstructure:"hosts"`
	Molly        MollyConfig        `mapstructure:"molly"`
	HeroPower    HeroPowerConfig    `mapstructure:"hero_power"`
	Sensitive    SensitiveConfig    `mapstructure:"sensitive"`
	AutoReply    AutoReplyConfig    `mapstructure:"auto_reply"`
	OnlineCourse OnlineCourseConfig `mapstructure:"online_course"`
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
	Enable      bool     `mapstructure:"enable"`
	Token       string   `mapstructure:"token"`
	Keywords    []string `mapstructure:"keywords"`
	AlertTimes  int      `mapstructure:"alert_times" json:"alert_times"`
	ShutSeconds int      `mapstructure:"shut_seconds" json:"shut_seconds"`
	Groups      []int64  `mapstructure:"groups"`
}

type AutoReplyConfig struct {
	Enable bool    `mapstructure:"enable"`
	Groups []int64 `mapstructure:"groups"`
}

type OnlineCourseConfig struct {
	Enable bool    `mapstructure:"enable"`
	Token  string  `mapstructure:"token"`
	Limit  int     `mapstructure:"limit"`
	Groups []int64 `mapstructure:"groups"`
}

func initConfig() Config {
	workDir, _ := os.Getwd()
	viper.AddConfigPath(workDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("application")
	initDefaultConfig()
	_ = viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Error("读取配置文件失败", err)
	}
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Error("解析结构体失败", err)
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
	sensitive.Keywords = make([]string, 0)
	data, _ = json.Marshal(*sensitive)
	_ = json.Unmarshal(data, &config)
	viper.SetDefault("sensitive", config)

	autoReply := new(AutoReplyConfig)
	autoReply.Enable = true
	viper.SetDefault("auto_reply", *autoReply)

	onlineCourse := new(OnlineCourseConfig)
	onlineCourse.Token = "free"
	onlineCourse.Enable = true
	onlineCourse.Limit = 1
	viper.SetDefault("online_course", *onlineCourse)

	autoReplyFile := "auto_reply_config.yaml"
	_, err := os.Stat(autoReplyFile)
	if os.IsNotExist(err) {
		AutoReplyList = append(AutoReplyList,
			map[string]interface{}{"ask": "你是谁", "answer": "我是机器人助手", "range": 0},
		)
		saveAutoReplyList(autoReplyFile)
	} else {
		loadAutoReplyList(autoReplyFile)
	}

	viper.SetDefault("hosts", []int64{})
	viper.SetDefault("name", "QQBotAssistant")
	viper.SetDefault("api_url", "http://127.0.0.1:8086")
}

func ReLoadSubConfig() {
	Molly = CONFIG.Molly
	HeroPower = CONFIG.HeroPower
	Sensitive = CONFIG.Sensitive
	AutoReply = CONFIG.AutoReply
	OnlineCourse = CONFIG.OnlineCourse
	saveAutoReplyList("auto_reply_config.yaml")
	loadAutoReplyList("auto_reply_config.yaml")
}

func saveAutoReplyList(fileName string) {
	data, err := yaml.Marshal(AutoReplyList)
	if err != nil {
		log.Error("序列化自动回复配置文件失败 Error:", err)
		return
	}
	if err = os.WriteFile(fileName, data, 0644); err != nil {
		log.Error("保存自动回复配置文件失败 Error:", err)
	}
}

func loadAutoReplyList(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Error("读取自动回复配置文件失败 Error:", err)
		return
	}
	if err = yaml.Unmarshal(data, &AutoReplyList); err != nil {
		log.Error("解析自动回复配置文件失败 Error:", err)
	}
}
