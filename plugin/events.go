package plugin

import (
	"QQBotAssistant/config"
	"QQBotAssistant/plugin/hero"
	"QQBotAssistant/plugin/molly"
	"QQBotAssistant/plugin/sensitive"
	"QQBotAssistant/util"
	"context"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
)

func LoadAllEvents(core *OPQBot.Core) {
	loadPluginSettingsEvent(core)

	if config.HeroPower.Enable {
		hero.LoadHeroPowerEvent(core)
	}
	if config.Molly.Enable {
		molly.LoadMollyEvent(core)
	}
	if config.Sensitive.Enable {
		sensitive.LoadSensitiveEvent(core)
	}
}

func loadPluginSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsHost(groupMsg.GetSenderUin()) || message == "" {
			return
		}
		if groupMsg.GetAtInfo() == nil {
			return
		}

		switch message {
		case config.HOST_ADD_KEY:
			for _, user := range groupMsg.GetAtInfo() {
				util.AddHost(user.Uin, "hosts")
			}
		case config.HOST_REMOVE_KEY:
			for _, user := range groupMsg.GetAtInfo() {
				util.DelHost(user.Uin, "hosts")
			}
		default:
			return
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, message+"成功")
	})
}
