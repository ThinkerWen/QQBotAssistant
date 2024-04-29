package reply

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"context"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/imroc/req/v3"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
	"net/url"
	"strconv"
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
		if !util.IsGroup(config.AutoReply.Groups, groupMsg.GetGroupUin()) || groupMsg.GetSenderUin() == event.GetCurrentQQ() {
			return
		}

		if answer := matchAsk(message, groupMsg.GetGroupUin()); answer != "" {
			var a *apiBuilder.Builder
			_ = json.Unmarshal([]byte(answer), &a)
			_ = DoAndResponse(a, event, ctx)
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
			rawMessage, _ := json.Marshal(parseSendMsg(event))
			currentState := checkReply(message, string(rawMessage), &v.State)
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

func parseSendMsg(event events.IEvent) (Builder *apiBuilder.Builder) {
	groupMsg := event.ParseGroupMsg()
	Builder = new(apiBuilder.Builder)
	CgiRequest := new(apiBuilder.CgiRequest)

	CgiCmd := "MessageSvc.PbSendMsg"
	Builder.CgiCmd = &CgiCmd
	Builder.CgiRequest = CgiRequest

	ToUin := groupMsg.GetGroupUin()
	CgiRequest.ToUin = &ToUin
	ToType := 2
	CgiRequest.ToType = &ToType

	Content := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
	CgiRequest.Content = &Content

	AtUinLists := make([]struct {
		Uin *int64 `json:"Uin,omitempty"`
	}, 0)
	for _, at := range groupMsg.GetAtInfo() {
		AtUinLists = append(AtUinLists, struct {
			Uin *int64 `json:"Uin,omitempty"`
		}{
			Uin: &at.Uin,
		})
	}
	CgiRequest.AtUinLists = AtUinLists

	Images := make([]*apiBuilder.File, 0)
	defer func() {
		if r := recover(); r != nil {
			log.Info(r)
		}
	}()
	for _, pic := range groupMsg.ParsePicMsg().GetPics() {
		file := new(apiBuilder.File)
		file.Url = pic.Url
		file.FileId = pic.FileId
		file.FileMd5 = pic.FileMd5
		file.FileSize = pic.FileSize
		Images = append(Images, file)
	}
	CgiRequest.Images = Images

	return
}

func DoAndResponse(b *apiBuilder.Builder, event events.IEvent, ctx context.Context) error {
	marshal, _ := json.Marshal(b)
	body := string(marshal)
	client := req.SetContext(ctx)
	u, _ := url.JoinPath(config.ApiUrl, "/v1/LuaApiCaller")
	client.SetURL(u)
	client.Method = "POST"
	resp := client.SetQueryParam("funcname", "MagicCgiCmd").SetQueryParam("qq", strconv.FormatInt(event.GetCurrentQQ(), 10)).SetBodyJsonString(body).Do()
	return resp.Err
}
