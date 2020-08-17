package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	host "github.com/team4yf/go-scheduler/pkg"
	"github.com/team4yf/go-scheduler/pkg/log"
)

var (
	hostname = host.GetHostname()
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

//SetUp 安装日志中间件
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		var requestData string
		contentType := c.ContentType()
		if c.Request.Method == "POST" {
			switch contentType {
			// 如果是 json 的，则将json取出来
			case gin.MIMEJSON:
				if c.Request.Body != nil {
					data, _ := c.GetRawData()
					// 这里需要将原来的数据还原回去，否则后面的handler获取不到原来的请求数据
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
					requestData = string(data)
				}
			}

		}

		//开始时间
		startTime := time.Now()

		//处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		//结束时间
		endTime := time.Now()

		//日志格式
		accessLogMap := make(map[string]interface{})

		accessLogMap["request_time"] = startTime
		accessLogMap["request_method"] = c.Request.Method
		accessLogMap["request_uri"] = c.Request.RequestURI
		accessLogMap["request_proto"] = c.Request.Proto
		accessLogMap["request_ua"] = c.Request.UserAgent()
		accessLogMap["request_content_type"] = contentType
		accessLogMap["request_referer"] = c.Request.Referer()
		accessLogMap["request_post_data"] = requestData
		accessLogMap["request_client_ip"] = c.ClientIP()

		accessLogMap["app_name"] = viper.GetString("name")
		accessLogMap["hostname"] = hostname

		accessLogMap["response_time"] = endTime
		accessLogMap["response_body"] = responseBody

		accessLogMap["cost_time"] = fmt.Sprintf("%v", endTime.Sub(startTime))

		accessLogJSON, _ := json.Marshal(accessLogMap)

		log.Infof((string)(accessLogJSON))
	}
}
