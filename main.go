package main

import (
	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/fpm-go-plugin-cron/plugin"
	_ "github.com/team4yf/fpm-go-plugin-email/plugin"
	_ "github.com/team4yf/fpm-go-plugin-orm/plugins/pg"
)

func main() {

	fpmApp := fpm.New()

	fpmApp.Init()

	fpmApp.Subscribe("#job/done", func(topic string, payload interface{}) {
		fpmApp.Logger.Debugf("topic: %s, payload: %v", topic, payload)
	})

	fpmApp.Run()
}
