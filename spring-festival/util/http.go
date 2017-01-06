package util

import (
	"crypto/tls"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

func HttpGetWithTimeout(requestUrl string, timeout time.Duration, encode string) string {
	if requestUrl == "" {
		return ""
	}
	req := httplib.Get(requestUrl)
	req.SetTimeout(timeout, timeout)
	str, err := req.String()
	if err != nil {
		beego.Error("http get url error", err)
		return ""
	}
	result, err := url.QueryUnescape(str)
	if err != nil {
		beego.Error("http get url error", err)
		return ""
	}
	return result
}

func HttpGet(url string) (result string) {
	if url == "" {
		return ""
	}
	timeout := time.Second * 2
	encode := "utf-8"
	return HttpGetWithTimeout(url, timeout, encode)
}

func HttpPost(requestUrl string, params map[string]string, headers map[string]string, timeout time.Duration) string {
	req := httplib.Post(requestUrl)
	if params != nil && len(params) > 0 {
		for key, value := range params {
			req.Param(key, value)
		}
	}
	if headers != nil && len(headers) > 0 {
		for key, value := range headers {
			req.Header(key, value)
		}
	}
	if strings.HasPrefix(requestUrl, "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	req.SetTimeout(timeout, timeout)
	str, err := req.String()
	if err != nil {
		beego.Error("http post url error", err)
		return ""
	}
	result, err := url.QueryUnescape(str)
	if err != nil {
		beego.Error("http post url error", err)
		return ""
	}
	return result
}
