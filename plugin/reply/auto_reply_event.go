package reply

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
)

var sequence = make([]Sequence, 0)

func LoadAutoReplyEvent(core *OPQBot.Core) {
	loadGroupEvent(core)
	loadSettingsEvent(core)
	log.Info("加载 自动回复助手 成功!")
}

func loadGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsGroup(config.AutoReply.Groups, groupMsg.GetGroupUin()) || groupMsg.GetSenderUin() == event.GetCurrentQQ() || message == "" {
			return
		}

		if answer := matchAsk(message, groupMsg.GetGroupUin()); answer != "" {
			_ = util.SendGroupMsg(event, groupMsg, ctx, answer)
			return
		}

		if message == config.REPLY_ADD {
			sequence = append(sequence, Sequence{Sender: groupMsg.GetSenderUin(), Receiver: groupMsg.GetGroupUin(), State: AutoReplyState{Step: 0}})
			_ = util.SendGroupMsg(event, groupMsg, ctx, config.REPLY_ASK)
			return
		}

		for i, v := range sequence {
			if !(v.Sender == groupMsg.GetSenderUin() && v.Receiver == groupMsg.GetGroupUin()) {
				continue
			}
			currentState := checkReply(message, &v.State)
			switch currentState {
			case 1:
				_ = util.SendGroupMsg(event, groupMsg, ctx, config.REPLY_ANSWER)
			case 2:
				_ = util.SendGroupMsg(event, groupMsg, ctx, config.REPLY_RANGE)
			case 3:
				if message == "1" {
					v.State.Range = groupMsg.GetGroupUin()
				}
				_ = util.SendGroupMsg(event, groupMsg, ctx, config.REPLY_ADD_SUCCESS)
				saveReply(v)
				sequence = removeSequence(groupMsg.GetSenderUin(), groupMsg.GetGroupUin())
			default:
				return
			}
			sequence[i] = v
			return
		}
	})
}

func loadSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsHost(groupMsg.GetSenderUin()) || message == "" {
			return
		}

		if config.REPLY_ON_KEY == message {
			util.AddGroup(config.AutoReply.Groups, groupMsg.GetGroupUin(), "auto_reply.groups")
		} else if config.REPLY_OFF_KEY == message {
			util.DelGroup(config.AutoReply.Groups, groupMsg.GetGroupUin(), "auto_reply.groups")
		} else {
			return
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, message+"成功")
	})
}
