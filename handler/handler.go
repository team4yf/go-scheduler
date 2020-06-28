package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/team4yf/go-scheduler/pkg/log"
	v "github.com/team4yf/go-scheduler/version"
)

type healthCheckResponse struct {
	Status   string    `json:"status"`
	Hostname string    `json:"hostname"`
	StartAt  time.Time `json:"startAt"`
	Version  string    `json:"version"`
	BuildAt  string    `json:"buildAt"`
}

var (
	// service  svc.LimitCouponActivityService
	// config   = conf.Config
	hostname = "undefined"
	startAt  = time.Now()
	version  = v.VERSION
	buildAt  = v.BuildAt
)

//HealthCheck 用于健康检查
func HealthCheck(c *gin.Context) {
	if c != nil {
		c.JSON(http.StatusOK, healthCheckResponse{Status: "UP", Hostname: hostname, StartAt: startAt, Version: version, BuildAt: buildAt})
		return
	}
	log.Infof("Not Health, gin.Context is nil")
}
