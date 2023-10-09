package request

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yicaoyimuys/GoGameServer/core/libs/logger"
	"go.uber.org/zap"
)

func init() {
	//设置http默认超时时间
	http.DefaultClient.Timeout = 3 * time.Second
	//不检测TLS证书
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// HttpGet
func HttpGet(url string, retryNum int) (string, uint32) {
	resp, err := http.Get(url)
	if err != nil {
		if retryNum > 0 {
			logger.Error("HttpGet 重试", zap.String("Url", url), zap.Error(err))
			return HttpGet(url, retryNum-1)
		} else {
			logger.Error("HttpGet", zap.String("Url", url), zap.Error(err))
			return "error", 1001
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("HttpGet", zap.String("Url", url), zap.Error(err))
		return "error", 1002
	}

	return string(body), 0
}

// HttpPost
func HttpPost(url string, retryNum int) (string, uint32) {
	arr := strings.Split(url, "?")
	resp, err := http.Post(arr[0], "application/x-www-form-urlencoded", strings.NewReader(arr[1]))
	if err != nil {
		if retryNum > 0 {
			logger.Error("HttpPost 重试", zap.String("Url", url), zap.Error(err))
			return HttpPost(url, retryNum-1)
		} else {
			logger.Error("HttpPost", zap.String("Url", url), zap.Error(err))
			return "error", 1001
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("HttpPost", zap.String("Url", url), zap.Error(err))
		return "error", 1002
	}

	return string(body), 0
}
