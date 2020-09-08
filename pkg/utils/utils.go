package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"

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
func SliceIndexOf(s []string, target string) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}
