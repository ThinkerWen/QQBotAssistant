package plugin

import (
	"QQBotAssistant/config"
	"github.com/opq-osc/OPQBot/v2"
)

func LoadAllEvents(core *OPQBot.Core) {
	if config.HeroPower.Enable {
		loadHeroPowerEvent(core)
	}
}
