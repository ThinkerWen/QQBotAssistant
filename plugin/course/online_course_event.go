package course

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
	"strings"
)

func LoadOnlineCourseEvent(core *OPQBot.Core) {
	loadGroupEvent(core)
	loadSettingsEvent(core)
	log.Info("加载 网课搜题助手 成功!")
}

func loadGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsGroup(config.OnlineCourse.Groups, groupMsg.GetGroupUin()) || message == "" {
			return
		}

		if message == config.COURSE_HELP_KEY {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.COURSE_HELP)
			return
		}
		params := strings.Split(message, " ")
		if config.COURSE_PFX != params[0] || len(params) < 2 {
			return
		}
		result := searchReason(params[1])
		if result == "" {
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.COURSE_NOT_FOUND)
			return
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, result)
	})
}

func loadSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsHost(groupMsg.GetSenderUin()) || message == "" {
			return
		}

		if config.COURSE_ON_KEY == message {
			util.AddGroup(config.OnlineCourse.Groups, groupMsg.GetGroupUin(), "online_course.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.COURSE_ON)
		} else if config.COURSE_OFF_KEY == message {
			util.DelGroup(config.OnlineCourse.Groups, groupMsg.GetGroupUin(), "online_course.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.COURSE_OFF)
		}
	})
}
