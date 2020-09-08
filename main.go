package main

import (
	"github.com/team4yf/go-scheduler/model"
	"github.com/team4yf/yf-fpm-server-go/fpm"

	_ "github.com/team4yf/fpm-go-plugin-email/plugin"
	_ "github.com/team4yf/fpm-go-plugin-orm/plugins/pg"
)

func main() {

	fpmApp := fpm.New()

	dbclient, _ := fpmApp.GetDatabase("pg")

	fpmApp.AddHook("BEFORE_INIT", func(_ *fpm.Fpm) {
		migrator := &model.Migration{
			DS: dbclient,
		}
		migrator.Install()
	}, 10)

	fpmApp.Init()

	fpmApp.AddBizModule("job", &fpm.BizModule{
		"execute": func(param *fpm.BizParam) (data interface{}, err error) {
			// start a job
			return
		},
	})

	fpmApp.Run()
}
