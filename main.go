package main

import (
	"QQBotAssistant/config"
	"QQBotAssistant/plugin"
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/spf13/viper"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	core, err := OPQBot.NewCore(config.ApiUrl, OPQBot.WithMaxRetryCount(5))
	if err != nil {
		log.Fatal(err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		_ = viper.Unmarshal(&config.CONFIG)
		config.ReLoadSubConfig()
	})

	plugin.LoadAllEvents(core)
	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
