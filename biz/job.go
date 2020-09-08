package biz

import (
	"github.com/team4yf/go-scheduler/model"
	"github.com/team4yf/go-scheduler/service/job"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

//InitJobBiz init the job biz
func InitJobBiz(fpmApp *fpm.Fpm) {
	dbclient, _ := fpmApp.GetDatabase("pg")
	jobService := job.NewSimpleJobService(dbclient)
	jobService.Init()
	if err := jobService.Start(); err != nil {
		fpmApp.Logger.Errorf("start scheduler error: %v", err)
	}

	fpmApp.AddBizModule("job", &fpm.BizModule{
		"execute": func(param *fpm.BizParam) (data interface{}, err error) {
			// start a job
			var job model.Job
			if err = param.Convert(&job); err != nil {
				return
			}
			jobService.Execute(&job, func(d interface{}, e error) {
				if e != nil {
					fpmApp.Logger.Errorf("execute error: %v", e)
					return
				}
				fpmApp.Logger.Infof("execute ok: %v", d)
			})
			return
		},
	})

	fpmApp.Subscribe("#job/success", func(topic string, data interface{}) {

	})
}
