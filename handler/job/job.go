package job

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/team4yf/go-scheduler/errno"
	"github.com/team4yf/go-scheduler/model"
	"github.com/team4yf/go-scheduler/pkg/log"
	"github.com/team4yf/go-scheduler/pkg/utils"
	"github.com/team4yf/go-scheduler/repository/job"

	jobSvc "github.com/team4yf/go-scheduler/service/job"
)

var (
	jobRepo    job.JobRepo
	jobService jobSvc.JobService
)

func Init() {
	jobRepo = job.NewJobRepo()

	jobService = jobSvc.NewSimpleJobService()
	jobService.Init()
	jobService.Start()
}

//List 获取全部的job信息
func List(c *gin.Context) {
	list, err := jobRepo.List(-1)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, list)
}

//Create create a job
func Create(c *gin.Context) {
	var req model.Job
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, errno.RequestBodyParseError, nil)
		return
	}
	err := jobRepo.Create(&req)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, 1)
}

func Get(c *gin.Context) {
	code := c.Param("code")
	job, err := jobRepo.Get(code)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, job)
}
func Execute(c *gin.Context) {
	code := c.Param("code")
	job, err := jobRepo.Get(code)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	var result interface{}
	var ex error
	var w sync.WaitGroup
	w.Add(1)
	jobService.Execute(job, func(data interface{}, err error) {
		result = data
		ex = err
		w.Done()
	})
	w.Wait()

	utils.SendResponse(c, ex, result)
}
func Update(c *gin.Context) {
	var req model.CommonMap
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	exists := false
	var status interface{}
	var code interface{}
	code, exists = req["code"]
	if !exists {
		utils.SendResponse(c, errors.New("Code required"), nil)
		return
	}
	status, exists = req["status"]
	if !exists {
		utils.SendResponse(c, errors.New("Status required"), nil)
		return
	}

	if err := jobRepo.Update(code.(string), req); err != nil {
		utils.SendResponse(c, err, nil)
		return
	}

	job := &model.Job{
		Code: code.(string),
	}
	switch status.(float64) {
	case 1:
		//run
		jobService.Restart(job, func(data interface{}, err error) {
			log.Errorf("data: %+v, err: %+v", data, err)
		})
	case 2:
		if err := jobService.Pause(job); err != nil {
			utils.SendResponse(c, err, nil)
			return
		}
	}

	utils.SendResponse(c, nil, 1)
}
