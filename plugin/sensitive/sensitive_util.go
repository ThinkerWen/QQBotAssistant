package sensitive

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/tidwall/gjson"
	"strings"
)

func isSensitive(content string) bool {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	data := fmt.Sprintf(`{"content":"%s","token":"%s"}`, content, config.Sensitive.Token)
	response, err := util.RequestPOST("https://www.hive-net.cn/funtools/sensitive/check", data, headers, nil)
	if err != nil || gjson.Get(string(response), "code").Int() != 0 {
		log.Error("敏感词检测请求异常 Error: ", err)
		return false
	}
	if gjson.Get(string(response), "data.minganCount").Int() != 0 {
		return true
	}
	return false
}

func isForbiddenKeyword(keyword string) bool {
	for _, k := range config.Sensitive.Keywords {
		if strings.Contains(keyword, k) {
			return true
		}
	}
	return false
}
