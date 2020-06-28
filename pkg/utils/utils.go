package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/team4yf/go-scheduler/errno"
	"github.com/teris-io/shortid"

	tnet "github.com/toolkits/net"
)

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

// GetLocalIP 获取本地内网IP
func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}

// GetBytes interface 转 byte
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Md5 字符串转md5
func Md5(str string) (string, error) {
	h := md5.New()

	_, err := io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	// 注意：这里不能使用string将[]byte转为字符串，否则会显示乱码
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//RespJSON the common json
type RespJSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SendResponse 返回json
func SendResponse(c *gin.Context, err error, data interface{}) {
	if err == nil {
		err = errno.OK
	}
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, RespJSON{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func ExportCsv(filePath string, header []string, data [][]string) (finalFilePath string, err error) {

	if finalFilePath, err = filepath.Abs(filePath); err != nil {
		return
	}

	dir := filepath.Dir(finalFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// not exists
		os.MkdirAll(dir, 0777)
		os.Chmod(dir, 0777)
	}

	file, err := os.Create(filePath)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write(header)
	for _, value := range data {
		err = writer.Write(value)
		if err != nil {
			return
		}
	}
	return
}

func SendFile(ctx *gin.Context, fileName, targetPath string) {
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(targetPath)
}

func Str2Int(str string, dft int) int {
	if str == "" {
		return dft
	}
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return (int)(i)
	}
	return dft

}

func JSON2String(j interface{}) (str string) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "{}"
	}
	str = (string)(bytes)
	return
}

func SliceIndexOf(s []string, target string) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

// GenShortID 生成一个id
func GenShortID() (string, error) {
	return shortid.Generate()
}

// GenUUID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}

// GetReqID 获取请求中的request_id
func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-ID")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}
