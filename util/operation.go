package util

import (
	"QQBotAssistant/config"
	"context"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
)

func SendGroupMsg(event events.IEvent, groupMsg events.IGroupMsg, ctx context.Context, message string) error {
	return apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg(message).Do(ctx)
}

func ShutGroupMember(event events.IEvent, groupMsg events.IGroupMsg, ctx context.Context, shutTime int) error {
	return apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).GroupManager().ProhibitedUser().ToGUin(groupMsg.GetGroupUin()).ToUid(groupMsg.GetSenderUid()).ShutTime(shutTime).Do(ctx)
}

func RevokeGroupMsg(event events.IEvent, groupMsg events.IGroupMsg, ctx context.Context) error {
	return apiBuilder.New(config.ApiUrl, event.GetCurrentQQ()).GroupManager().RevokeMsg().ToGUin(groupMsg.GetGroupUin()).MsgSeq(groupMsg.GetMsgSeq()).MsgRandom(groupMsg.GetMsgRandom()).Do(ctx)
}
