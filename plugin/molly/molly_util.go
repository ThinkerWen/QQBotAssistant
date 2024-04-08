package molly

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"encoding/json"
	"log"
)

type ContentMolly struct {
	Content  string `json:"content"`
	Type     int    `json:"type"`
	From     string `json:"from"`
	FromName string `json:"fromName"`
	To       string `json:"to"`
	ToName   string `json:"toName"`
}

func mollyChat(content ContentMolly) string {
	headers := make(map[string]string)
	headers["Api-Key"] = config.Molly.ApiKey
	headers["Api-Secret"] = config.Molly.ApiSecret
	headers["Content-Type"] = "application/json;charset=UTF-8"
	if data, errJson := json.Marshal(content); errJson == nil {
		if response, err := util.RequestPOST("https://api.mlyai.com/reply", string(data), headers, nil); err == nil {
			return string(response)
		}
	}
	log.Println("Molly 聊天调用失败")
	return ""
}
