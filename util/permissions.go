package util

import (
	"QQBotAssistant/config"
	"github.com/spf13/viper"
)

func IsHost(host int64) bool {
	for _, v := range config.CONFIG.Hosts {
		if v == host {
			return true
		}
	}
	return false
}

func IsGroup(groups []int64, group int64) bool {
	for _, v := range groups {
		if v == group {
			return true
		}
	}
	return false
}

func AddHost(host int64, key string) {
	if IsHost(host) {
		return
	}
	config.CONFIG.Hosts = append(config.CONFIG.Hosts, host)
	viper.Set(key, config.CONFIG.Hosts)
	_ = viper.WriteConfig()
}

func DelHost(host int64, key string) {
	var result []int64
	for _, v := range config.CONFIG.Hosts {
		if v != host {
			result = append(result, v)
		}
	}
	viper.Set(key, result)
	_ = viper.WriteConfig()
}

func AddGroup(groups []int64, group int64, key string) {
	if IsGroup(groups, group) {
		return
	}
	groups = append(groups, group)
	viper.Set(key, groups)
	_ = viper.WriteConfig()
}

func DelGroup(groups []int64, group int64, key string) {
	var result []int64
	for _, v := range groups {
		if v != group {
			result = append(result, v)
		}
	}
	viper.Set(key, result)
	_ = viper.WriteConfig()
}
