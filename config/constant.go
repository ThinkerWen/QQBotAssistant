package config

const (
	PFX             = "查战力"
	HELP_KEY        = "查战力帮助"
	HOST_ADD_KEY    = "添加插件管理员"
	HOST_REMOVE_KEY = "移除插件管理员"

	// ================================ HeroPower ================================

	HERO_ON_KEY       = "开启战力查询"
	HERO_OFF_KEY      = "关闭战力查询"
	HERO_ON           = "已开启战力查询！\n输入“查战力帮助”，获取战力查询方法"
	HERO_OFF          = "已关闭战力查询，需要开启请找管理员"
	HERO_WRONG_TOKEN  = "查询指令错误！\n输入“查战力帮助”，获取战力查询方法"
	HERO_HELP         = "查战力请输入：\n“查战力 英雄 区服”\n例如：\n\n查战力 李白 安卓QQ\n查战力 李白 苹果微信\n\n不要漏掉空格"
	HERO_POWER_RESULT = "查询结果如下：\n\n系统：%s\n英雄：%s\n\n更新时间：\n%s\n\n省标：\n%s %s分\n市标：\n%s %s分\n区标：\n%s %s分\n\n微信小程序《峡谷战力查改》"

	// ================================ Molly ================================

	MOLLY_ON_KEY  = "开启聊天机器人"
	MOLLY_OFF_KEY = "关闭聊天机器人"
	MOLLY_ON      = "%s已经回来啦\n请@我来和我聊天哦~"
	MOLLY_OFF     = "%s走了，有机会再见～"

	// ================================ Sensitive ================================

	SENSITIVE_ON_KEY  = "开启敏感词过滤"
	SENSITIVE_OFF_KEY = "关闭敏感词过滤"
	SENSITIVE_ON      = "本群已开启24h监控，请注意言行"
	SENSITIVE_OFF     = "监控不到大家了，想干什么GKD！"
)
