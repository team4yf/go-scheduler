//Package task
package task

import (
	"github.com/gin-gonic/gin"

	"github.com/team4yf/go-scheduler/handler/task"
)

//Load 加载task相关的接口
func Load(g *gin.RouterGroup, mw ...gin.HandlerFunc) {
	task.Init()
	group := g.Group("task", mw...)
	group.GET("/list/:code", task.List)
	group.GET("/detail/:id", task.Get)
	group.GET("/export/:code", task.Export)
}
