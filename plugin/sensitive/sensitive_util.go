package sensitive

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
)

func isSensitive(content string) bool {
	headers := make(map[string]string)
	headers["Connection"] = "keep-alive"
	headers["Accept-Language"] = "zh-CN,zh;q=0.9"
	headers["Origin"] = "http://www.zhipaiwu.com"
	headers["X-Requested-With"] = "XMLHttpRequest"
	headers["Accept"] = "application/json, text/javascript, */*; q=0.01"
	headers["Referer"] = "http://www.zhipaiwu.com/index.php/Weijinci/index.html"
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"

	data := fmt.Sprintf(`content=%s`, url.QueryEscape(content))
	response, err := util.RequestPOST("http://www.zhipaiwu.com/index.php/Weijinci/postIndex.html", data, headers, nil)
	if err != nil || gjson.Get(string(response), "code").Int() != 200 {
		log.Error("敏感词检测请求异常 Error: ", err)
		return false
	}
	if gjson.Get(string(response), "result.minganCount").Int() != 0 {
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
