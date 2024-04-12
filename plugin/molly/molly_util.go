package molly

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
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
	client := resty.New()
	client.SetHeaders(headers)
	client.SetTimeout(10 * time.Second)
	if data, errJson := json.Marshal(content); errJson == nil {
		if response, err := util.RequestPOST("https://api.mlyai.com/reply", string(data), headers, client); err == nil {
			return string(response)
		}
	}
	log.Println("Molly 聊天调用失败")
	return ""
}
