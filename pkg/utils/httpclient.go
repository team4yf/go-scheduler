package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

type HTTPAuthType string
type HTTPAuthData map[string]interface{}

const (
	TypeHTTPAuthBasic = "basic"
)

type HttpAuth struct {
	Type HTTPAuthType
	Data HTTPAuthData
}

func Get(url string, timeout int) ResponseWrapper {
	return GetWithAuth(url, timeout, nil)
}
func PostJson(url string, body string, timeout int) ResponseWrapper {
	return PostJsonWithAuth(url, body, timeout, nil)
}
func PostForm(url string, body string, timeout int) ResponseWrapper {
	return PostFormWithAuth(url, body, timeout, nil)
}

func PostFormWithAuth(url string, params string, timeout int, auth *HttpAuth) ResponseWrapper {
	return PostWithAuth(url, "application/x-www-form-urlencoded", params, timeout, auth)
}

func PostJsonWithAuth(url string, body string, timeout int, auth *HttpAuth) ResponseWrapper {
	return PostWithAuth(url, "application/json", body, timeout, auth)
}

func GetWithAuth(url string, timeout int, auth *HttpAuth) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}

	return request(req, timeout, auth)
}

func PostWithAuth(url string, contentType string, params string, timeout int, auth *HttpAuth) ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return request(req, timeout, auth)
}

func request(req *http.Request, timeout int, auth *HttpAuth) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req)
	if auth != nil {
		switch auth.Type {
		case TypeHTTPAuthBasic:
			req.SetBasicAuth(auth.Data["username"].(string), auth.Data["password"].(string))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header

	return wrapper
}

func setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "golang/gocron")
}

func createRequestError(err error) ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return ResponseWrapper{0, errorMessage, make(http.Header)}
}
