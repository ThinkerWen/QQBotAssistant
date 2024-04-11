package course

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"time"
)

func searchReason(param string) string {
	res := ""
	link := fmt.Sprintf("https://www.hive-net.cn/backend/wangke/search?token=%s&question=%s", config.OnlineCourse.Token, param)
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	response, err := util.RequestGET(link, nil, client)
	if err != nil || gjson.Get(string(response), "code").Int() != 0 {
		log.Error("搜索接口调用失败 Error: ", err)
		return ""
	}
	result := string(response)
	for i, answer := range gjson.Get(result, "data.reasonList").Array() {
		question := answer.Get("question").Str
		options := answer.Get("options").Str
		reason := answer.Get("reason").Str
		data := fmt.Sprintf("================\n问题:\n%s", question)
		if options != "" && options != "无" {
			data += fmt.Sprintf("\n选项:\n%s", options)
		}
		data += fmt.Sprintf("\n答案:\n%s\n", reason)
		res += data
		if i == config.OnlineCourse.Limit-1 {
			return res
		}
	}
	return res
}
