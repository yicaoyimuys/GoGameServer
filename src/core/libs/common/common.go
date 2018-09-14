package common

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

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
