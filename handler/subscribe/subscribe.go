package subscribe

import (
	"github.com/gin-gonic/gin"
	"github.com/team4yf/go-scheduler/errno"
	"github.com/team4yf/go-scheduler/model"
	"github.com/team4yf/go-scheduler/pkg/utils"
	"github.com/team4yf/go-scheduler/repository/subscribe"
)

var subRep subscribe.SubscribeRepo

type subObj struct {
	NotifyType  string `json:"notifyType"`
	Subscriber  string `json:"subscriber"`
	NotifyEvent string `json:"notifyEvent"`
}

type subReq struct {
	Data []subObj `json:"data"`
}

func Init() {
	subRep = subscribe.NewSubscribeRepo()
}

func Sub(c *gin.Context) {
	code := c.Param("code")
	var req subReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, errno.RequestBodyParseError, nil)
		return
	}
	var subModels []*model.Subscribe
	for _, sub := range req.Data {
		subModel := &model.Subscribe{
			Code:        code,
			Subscriber:  sub.Subscriber,
			NotifyEvent: sub.NotifyEvent,
			NotifyType:  sub.NotifyType,
			NotifyTopic: code + "/" + sub.NotifyEvent,
		}
		subModels = append(subModels, subModel)
	}
	if err := subRep.Subscribe(subModels); err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, 1)

}

func UnSub(c *gin.Context) {
	code := c.Param("code")
	var req model.Subscribe
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, errno.RequestBodyParseError, nil)
		return
	}
	req.Code = code
	if err := subRep.UnSubscribe(&req); err != nil {
		utils.SendResponse(c, err, nil)
		return
	}
	utils.SendResponse(c, nil, 1)
}

func List(c *gin.Context) {
	code := c.Param("code")

	list, err := subRep.List(code)
	if err != nil {
		utils.SendResponse(c, err, nil)
		return
	}

	utils.SendResponse(c, nil, list)

}
