package molly

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/events"
	"github.com/opq-osc/OPQBot/v2/faces"
	"github.com/tidwall/gjson"
	"strconv"
)

func LoadMollyEvent(core *OPQBot.Core) {
	loadGroupEvent(core)
	loadSettingsEvent(core)
	log.Info("加载 莫莉云机器人 成功!")
}

func loadGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsGroup(config.Molly.Groups, groupMsg.GetGroupUin()) || message == "" || groupMsg.GetAtInfo() == nil {
			return
		}
		if groupMsg.GetAtInfo()[0].Uin != config.Molly.QQ {
			return
		}
		data := new(ContentMolly)
		data.Type = 2
		data.Content = message
		data.ToName = groupMsg.GetGroupInfo().GroupName
		data.FromName = groupMsg.GetSenderNick()
		data.To = strconv.FormatInt(groupMsg.GetGroupUin(), 10)
		data.From = strconv.FormatInt(groupMsg.GetSenderUin(), 10)
		result := mollyChat(*data)
		if result == "" || gjson.Get(result, "code").Str != "00000" || len(gjson.Get(result, "data").Array()) == 0 {
			_ = util.SendGroupMsg(event, groupMsg, ctx, "我现在不想说话"+faces.Face_doge)
			return
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, gjson.Get(result, "data|0.content").Str)
	})
}

func loadSettingsEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if !util.IsHost(groupMsg.GetSenderUin()) || message == "" {
			return
		}

		if config.MOLLY_ON_KEY == message {
			util.AddGroup(config.Molly.Groups, groupMsg.GetGroupUin(), "molly.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, fmt.Sprintf(config.MOLLY_ON, config.Molly.Name))
		} else if config.MOLLY_OFF_KEY == message {
			util.DelGroup(config.Molly.Groups, groupMsg.GetGroupUin(), "molly.groups")
			_ = util.SendGroupMsg(event, groupMsg, ctx, fmt.Sprintf(config.MOLLY_OFF, config.Molly.Name))
		}
	})
}
