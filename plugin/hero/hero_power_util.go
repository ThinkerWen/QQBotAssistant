package hero

import (
	"QQBotAssistant/config"
	"QQBotAssistant/util"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

func getHeroPower(token, hero, server string) string {
	api := fmt.Sprintf("https://www.hive-net.cn/funtools/heroPower/getPower?hero=%s&type=%s&token=%s", hero, server, token)
	response, err := util.RequestGET(api, nil, nil)
	if err != nil || gjson.Get(string(response), "code").Int() != 0 {
		return "请求出错，请联系作者！"
	}
	data := gjson.Get(string(response), "data").String()
	return fmt.Sprintf(config.HERO_POWER_RESULT, gjson.Get(data, "server").Str, gjson.Get(data, "name").Str, gjson.Get(data, "updatetime").Str,
		gjson.Get(data, "province.name").Str, gjson.Get(data, "province.power").Str, gjson.Get(data, "city.name").Str,
		gjson.Get(data, "city.power").Str, gjson.Get(data, "area.name").Str, gjson.Get(data, "area.power").Str)
}

func getHeroServer(server string) string {
	switch strings.ToUpper(server) {
	case "安卓QQ":
		return "aqq"
	case "安卓微信":
		return "awx"
	case "苹果QQ":
		return "ios_qq"
	case "苹果微信":
		return "ios_wx"
	default:
		return ""
	}
}

func isHeroHost(host int64) bool {
	return util.IsHost(config.HeroPower.Hosts, host)
}

func isHeroGroup(group int64) bool {
	return util.IsGroup(config.HeroPower.Groups, group)
}

func addHeroHost(host int64) {
	util.AddHost(config.HeroPower.Hosts, host, "hero_power.hosts")
}

func delHeroHost(host int64) {
	util.DelHost(config.HeroPower.Hosts, host, "hero_power.hosts")
}

func addHeroGroup(group int64) {
	util.AddGroup(config.HeroPower.Groups, group, "hero_power.groups")
}

func delHeroGroup(group int64) {
	util.DelGroup(config.HeroPower.Groups, group, "hero_power.groups")
}