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
	loadGroupEvent(core)
	loadSettingsEvent(core)
	log.Info("加载 王者荣耀战力查询 成功!")
}

func loadGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsGroup(config.HeroPower.Groups, groupMsg.GetGroupUin()) || message == "" {
			return
		}

		if message == config.HERO_HELP_KEY {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HERO_HELP)
			return
		}
		params := strings.Split(message, " ")
		if config.HERO_PFX != params[0] {
			return
		}
		if len(params) != 3 {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HERO_WRONG_TOKEN)
			return
		}

		hero := params[1]
		server := getHeroServer(params[2])
		if hero == "" || server == "" {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HERO_WRONG_TOKEN)
			return
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, getHeroPower(config.HeroPower.Token, hero, server))
	})
}

func loadSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsHost(groupMsg.GetSenderUin()) || message == "" {
			return
		}

		if config.HERO_ON_KEY == message {
			util.AddGroup(config.HeroPower.Groups, groupMsg.GetGroupUin(), "hero_power.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HERO_ON)
		} else if config.HERO_OFF_KEY == message {
			util.AddGroup(config.HeroPower.Groups, groupMsg.GetGroupUin(), "hero_power.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.HERO_OFF)
		}
	})
}
