package sensitive

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"github.com/charmbracelet/log"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"strings"
)

var sensitiveCount = make(map[int64]int)

func LoadSensitiveEvent(core *OPQBot.Core) {
	loadGroupEvent(core)
	log.Info("加载 Sensitive 成功!")
}

func loadGroupEvent(core *OPQBot.Core) {
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		groupMsg := event.ParseGroupMsg()
		message := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
		if strings.TrimSpace(message) == "" || groupMsg.GetSenderUin() == event.GetCurrentQQ() {
			return
		}
		if !isSensitive(message) {
			return
		}
		if _, ok := sensitiveCount[groupMsg.GetSenderUin()]; !ok {
			sensitiveCount[groupMsg.GetSenderUin()] = 1
		} else if sensitiveCount[groupMsg.GetSenderUin()] < 2 {
			sensitiveCount[groupMsg.GetSenderUin()]++
		} else {
			delete(sensitiveCount, groupMsg.GetSenderUin())
			_ = apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).GroupManager().ProhibitedUser().ToGUin(groupMsg.GetGroupUin()).ToUid(groupMsg.GetSenderUid()).ShutTime(60).Do(ctx)
		}
		_ = util.SendGroupMsg(event, groupMsg, ctx, "请勿发送不当言论，达到3次将禁言")
		_ = apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).GroupManager().RevokeMsg().ToGUin(groupMsg.GetGroupUin()).MsgSeq(groupMsg.GetMsgSeq()).MsgRandom(groupMsg.GetMsgRandom()).Do(ctx)
	})
}
