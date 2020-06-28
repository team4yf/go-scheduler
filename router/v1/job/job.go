//Package job the routers for the job
package job

import (
	"github.com/gin-gonic/gin"

	"github.com/team4yf/go-scheduler/handler/job"
)

//Load 加载Store相关的接口
func Load(g *gin.RouterGroup, mw ...gin.HandlerFunc) {
	job.Init()
	group := g.Group("job", mw...)
	group.GET("/list", job.List)
	group.GET("/get/:code", job.Get)
	group.GET("/execute/:code", job.Execute)
	group.POST("/create", job.Create)
	group.POST("/update", job.Update)
}
