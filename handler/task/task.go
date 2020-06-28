package task

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/team4yf/go-scheduler/config"
	"github.com/team4yf/go-scheduler/pkg/utils"
	"github.com/team4yf/go-scheduler/repository/task"
)

var (
	taskRepo task.TaskRepo
	headers  []string
	keys     []string
)

func Init() {
	taskRepo = task.NewTaskRepo()
	headers = strings.Split("id,任务编号,执行结果,开始时间,结束时间,执行状态,耗时(毫秒),URL", ",")
	keys = strings.Split("ID,Code,Log,StartAt,EndAt,Status,Cost,URL", ",")
}

//List 获取全部的task信息
func List(c *gin.Context) {
	code := c.Param("code")
	p := utils.Str2Int(c.Query("p"), 1)

	l := utils.Str2Int(c.Query("l"), 10)

	list, err := taskRepo.List(code, p, l)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, list)
}

//Export 获取全部的task信息
func Export(c *gin.Context) {
	code := c.Param("code")
	p := utils.Str2Int(c.Query("p"), 1)

	l := utils.Str2Int(c.Query("l"), 10)

	list, err := taskRepo.List(code, p, l)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	if len(list) < 1 {
		utils.SendResponse(c, errors.New("no data"), nil)
		return
	}
	data := [][]string{}
	for _, one := range list {
		row := one.String(keys)
		data = append(data, row)
	}
	finalFilePath, err := utils.ExportCsv(filepath.Join(config.GetString("export", "./export"), "data.csv"), headers, data)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendFile(c, "data.csv", finalFilePath)
}

func Get(c *gin.Context) {
	id := c.Param("id")
	u, _ := strconv.ParseUint(id, 10, 32)

	task, err := taskRepo.Get((uint)(u))
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, task)
}
