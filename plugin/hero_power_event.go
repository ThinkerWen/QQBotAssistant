package plugin

import (
	"QQBotAssistant/config"
	"context"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/spf13/viper"
	"strings"
)

func loadHeroPowerEvent(core *OPQBot.Core) {
	loadHeroPowerGroupEvent(core)
	loadHeroPowerSettingsEvent(core)
	log.Info("加载 HeroPower 成功!")
}

func loadHeroPowerGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !isHeroGroup(groupMsg.GetGroupUin()) || strings.TrimSpace(message) == "" {
			return
		}

		if message == config.HELP_KEY {
			_ = sendMsg(event, groupMsg, ctx, config.HELP)
			return
		}
		params := strings.Split(message, " ")
		if config.PFX != params[0] {
			return
		}
		if len(params) != 3 {
			_ = sendMsg(event, groupMsg, ctx, config.WRONG_TOKEN)
			return
		}

		hero := params[1]
		server := getHeroServer(params[2])
		if strings.TrimSpace(hero) == "" || strings.TrimSpace(server) == "" {
			_ = sendMsg(event, groupMsg, ctx, config.WRONG_TOKEN)
			return
		}
		_ = sendMsg(event, groupMsg, ctx, getHeroPower(config.HeroPower.Token, hero, server))
	})
}

func loadHeroPowerSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !isHeroHost(groupMsg.GetSenderUin()) || strings.TrimSpace(message) == "" {
			return
		}

		params := strings.Split(message, " ")
		if len(params[0]) == 2 && groupMsg.GetAtInfo() == nil {
			return
		}
		if config.SEARCH_ON_KEY == message {
			addHeroGroup(groupMsg.GetGroupUin())
			_ = sendMsg(event, groupMsg, ctx, config.SEARCH_ON)
		} else if config.SEARCH_OFF_KEY == message {
			delHeroGroup(groupMsg.GetGroupUin())
			_ = sendMsg(event, groupMsg, ctx, config.SEARCH_OFF)
		} else if config.HOST_ADD_KEY == params[0] {
			addHeroHost(groupMsg.GetAtInfo()[0].Uin)
			_ = sendMsg(event, groupMsg, ctx, config.HOST_ADD_KEY+"成功")
		} else if config.HOST_REMOVE_KEY == params[0] {
			delHeroHost(groupMsg.GetAtInfo()[0].Uin)
			_ = sendMsg(event, groupMsg, ctx, config.HOST_REMOVE_KEY+"成功")
		}
	})
}

func isHeroHost(host int64) bool {
	for _, v := range config.HeroPower.Hosts {
		if v == host {
			return true
		}
	}
	return false
}

func isHeroGroup(group int64) bool {
	for _, v := range config.HeroPower.Groups {
		if v == group {
			return true
		}
	}
	return false
}

func addHeroHost(host int64) {
	if isHeroHost(host) {
		return
	}
	config.HeroPower.Hosts = append(config.HeroPower.Hosts, host)
	viper.Set("hero_power.hosts", config.HeroPower.Hosts)
	_ = viper.WriteConfig()
}

func delHeroHost(host int64) {
	var result []int64
	for _, v := range config.CONFIG.HeroPower.Hosts {
		if v != host {
			result = append(result, v)
		}
	}
	viper.Set("hero_power.hosts", result)
	_ = viper.WriteConfig()
}

func addHeroGroup(group int64) {
	if isHeroGroup(group) {
		return
	}
	config.HeroPower.Groups = append(config.HeroPower.Groups, group)
	viper.Set("hero_power.groups", config.HeroPower.Groups)
	_ = viper.WriteConfig()
}

func delHeroGroup(group int64) {
	var result []int64
	for _, v := range config.CONFIG.HeroPower.Groups {
		if v != group {
			result = append(result, v)
		}
	}
	viper.Set("hero_power.groups", result)
	_ = viper.WriteConfig()
}

func sendMsg(event events.IEvent, groupMsg events.IGroupMsg, ctx context.Context, message string) error {
	return apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg(message).Do(ctx)
}
