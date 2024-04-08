package hero

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
	"strings"
)

func LoadHeroPowerEvent(core *OPQBot.Core) {
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
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HELP)
			return
		}
		params := strings.Split(message, " ")
		if config.PFX != params[0] {
			return
		}
		if len(params) != 3 {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.WRONG_TOKEN)
			return
		}

		hero := params[1]
		server := getHeroServer(params[2])
		if strings.TrimSpace(hero) == "" || strings.TrimSpace(server) == "" {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.WRONG_TOKEN)
			return
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, getHeroPower(config.HeroPower.Token, hero, server))
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
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.SEARCH_ON)
		} else if config.SEARCH_OFF_KEY == message {
			delHeroGroup(groupMsg.GetGroupUin())
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.SEARCH_OFF)
		} else if config.HOST_ADD_KEY == params[0] {
			addHeroHost(groupMsg.GetAtInfo()[0].Uin)
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HOST_ADD_KEY+"成功")
		} else if config.HOST_REMOVE_KEY == params[0] {
			delHeroHost(groupMsg.GetAtInfo()[0].Uin)
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HOST_REMOVE_KEY+"成功")
		}
	})
}
