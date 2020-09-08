package main

import (
	"github.com/team4yf/go-scheduler/biz"
	"github.com/team4yf/go-scheduler/model"
	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/fpm-go-plugin-email/plugin"
	_ "github.com/team4yf/fpm-go-plugin-orm/plugins/pg"
)

func main() {

	fpmApp := fpm.New()

	fpmApp.AddHook("BEFORE_INIT", func(_ *fpm.Fpm) {
		dbclient, _ := fpmApp.GetDatabase("pg")
		migrator := &model.Migration{
			DS: dbclient,
		}
		migrator.Install()
	}, 10)

	fpmApp.Init()

	biz.InitJobBiz(fpmApp)

	fpmApp.Run()
}
