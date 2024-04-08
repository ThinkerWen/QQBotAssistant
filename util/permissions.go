package util

import (
	"QQBotAssistant/config"
	"context"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/spf13/viper"
)

func IsHost(hosts []int64, host int64) bool {
	for _, v := range hosts {
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

func AddHost(hosts []int64, host int64, key string) {
	if IsHost(hosts, host) {
		return
	}
	hosts = append(hosts, host)
	viper.Set("hero_power.hosts", hosts)
	_ = viper.WriteConfig()
}

func DelHost(hosts []int64, host int64, key string) {
	var result []int64
	for _, v := range hosts {
		if v != host {
			result = append(result, v)
		}
	}
	viper.Set("hero_power.hosts", result)
	_ = viper.WriteConfig()
}

func AddGroup(groups []int64, group int64, key string) {
	if IsGroup(groups, group) {
		return
	}
	groups = append(groups, group)
	viper.Set("hero_power.groups", groups)
	_ = viper.WriteConfig()
}

func DelGroup(groups []int64, group int64, key string) {
	var result []int64
	for _, v := range groups {
		if v != group {
			result = append(result, v)
		}
	}
	viper.Set("hero_power.groups", result)
	_ = viper.WriteConfig()
}

func SendGroupMsg(event events.IEvent, groupMsg events.IGroupMsg, ctx context.Context, message string) error {
	return apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg(message).Do(ctx)
}
