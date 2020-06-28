package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/team4yf/go-scheduler/config"
	"github.com/team4yf/go-scheduler/handler"
	"github.com/team4yf/go-scheduler/pkg/email"
	"github.com/team4yf/go-scheduler/pkg/log"
	"github.com/team4yf/go-scheduler/router/v1/job"
	"github.com/team4yf/go-scheduler/router/v1/subscribe"
	"github.com/team4yf/go-scheduler/router/v1/task"

	"github.com/team4yf/go-scheduler/middleware"
	"github.com/team4yf/go-scheduler/model"
)

var migration model.Migration

func main() {
	config.Init("")
	// Set gin mode.
	gin.SetMode(viper.GetString("mode"))

	// Init the model
	model.CreateDb()
	migration.Install()

	// Init email
	email.Init()
	// Create the Gin engine.
	engine := gin.Default()

	// HealthCheck 健康检查路由
	engine.GET("/health", handler.HealthCheck)
	// metrics router 可以在 prometheus 中进行监控
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	loggerMdl := middleware.SetUp()
	// API Routes.
	group := engine.Group("api/v1", loggerMdl)
	{
		job.Load(group)
		task.Load(group)
		subscribe.Load(group)
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	srv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()
	gracefulStop(srv)
}

// gracefulStop 优雅退出
// 等待中断信号以超时 5 秒正常关闭服务器
// 官方说明：https://github.com/gin-gonic/gin#graceful-restart-or-stop
func gracefulStop(srv *http.Server) {
	quit := make(chan os.Signal)
	// kill 命令发送信号 syscall.SIGTERM
	// kill -2 命令发送信号 syscall.SIGINT
	// kill -9 命令发送信号 syscall.SIGKILL
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// 5 秒后捕获 ctx.Done() 信号
	select {
	case <-ctx.Done():
		log.Info("timeout of 5 seconds.")
	default:
	}
	log.Info("Server exiting")
}
