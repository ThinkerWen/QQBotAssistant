package plugin

import (
	"QQBotAssistant/config"
	"QQBotAssistant/plugin/hero"
	"QQBotAssistant/plugin/molly"
	"github.com/opq-osc/OPQBot/v2"
)

func LoadAllEvents(core *OPQBot.Core) {
	if config.HeroPower.Enable {
		hero.LoadHeroPowerEvent(core)
	}
	if config.Molly.Enable {
		molly.LoadMollyEvent(core)
	}
}
