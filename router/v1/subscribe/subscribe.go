//Package subscribe
package subscribe

import (
	"github.com/gin-gonic/gin"

	"github.com/team4yf/go-scheduler/handler/subscribe"
)

//Load 加载subscribe相关的接口
func Load(g *gin.RouterGroup, mw ...gin.HandlerFunc) {
	subscribe.Init()
	group := g.Group("subscribe", mw...)
	group.POST("/sub/:code", subscribe.Sub)
	group.POST("/unSub/:code", subscribe.UnSub)
	group.GET("/list/:code", subscribe.List)
}
