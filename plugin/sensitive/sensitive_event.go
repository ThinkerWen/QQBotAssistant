package sensitive

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/spf13/viper"
	"strings"
)

var sensitiveCount = make(map[int64]int)

func LoadSensitiveEvent(core *OPQBot.Core) {
	loadGroupEvent(core)
	loadSettingsEvent(core)
	log.Info("加载 敏感词检测 成功!")
}

func loadGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsGroup(config.Sensitive.Groups, groupMsg.GetGroupUin()) || message == "" || groupMsg.GetSenderUin() == event.GetCurrentQQ() {
			return
		}
		if !isForbiddenKeyword(message) && !isSensitive(message) {
			return
		}

		if _, ok := sensitiveCount[groupMsg.GetSenderUin()]; !ok {
			sensitiveCount[groupMsg.GetSenderUin()] = 1
		} else if sensitiveCount[groupMsg.GetSenderUin()] < config.Sensitive.AlertTimes-1 {
			sensitiveCount[groupMsg.GetSenderUin()]++
		} else {
			delete(sensitiveCount, groupMsg.GetSenderUin())
			_ = util.ShutGroupMember(event, groupMsg, ctx, config.Sensitive.ShutSeconds)
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, fmt.Sprintf("请勿发送不当言论，达到%d次将禁言", config.Sensitive.AlertTimes))
		_ = util.RevokeGroupMsg(event, groupMsg, ctx)
	})
}

func loadSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsHost(groupMsg.GetSenderUin()) || message == "" {
			return
		}

		if config.SENSITIVE_ON_KEY == message {
			util.AddGroup(config.Sensitive.Groups, groupMsg.GetGroupUin(), "sensitive.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.SENSITIVE_ON)
		} else if config.SENSITIVE_OFF_KEY == message {
			util.AddGroup(config.Sensitive.Groups, groupMsg.GetGroupUin(), "sensitive.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.SENSITIVE_OFF)
		} else if strings.Contains(message, config.SENSITIVE_ADD_KEY) {
			params := strings.Split(message, " ")
			if len(params) != 2 {
				return
			}
			config.Sensitive.Keywords = append(config.Sensitive.Keywords, params[1])
			viper.Set("sensitive.keywords", config.Sensitive.Keywords)
			_ = viper.WriteConfig()
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.SENSITIVE_ADD_KEY+"成功")
		}
	})
}
