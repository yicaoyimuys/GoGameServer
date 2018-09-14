package common

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	. "tools"
)

func init() {
	//设置http默认超时时间
	http.DefaultClient.Timeout = 3 * time.Second
	//不检测TLS证书
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

//HttpGet
func HttpGet(url string) (string, uint32) {
	resp, err := http.Get(url)
	if err != nil {
		ERR("HttpGet", url, err)
		return "error", 1001
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERR("HttpGet", url, err)
		return "error", 1002
	}

	return string(body), 0
}

//HttpPost
func HttpPost(url string, retryNum int) (string, uint32) {
	arr := strings.Split(url, "?")
	resp, err := http.Post(arr[0], "application/x-www-form-urlencoded", strings.NewReader(arr[1]))
	if err != nil {
		if retryNum > 0 {
			ERR("HttpPost 重试", url, err)
			return HttpPost(url, retryNum-1)
		} else {
			ERR("HttpPost", url, err)
			return "error", 1001
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERR("HttpPost", url, err)
		return "error", 1002
	}

	return string(body), 0
}

//Md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)
	return md5Str
}

//获取当前毫秒时间戳
func UnixMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
