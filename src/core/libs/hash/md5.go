package hash

import (
	"crypto/md5"
	"encoding/hex"
)

//Md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	md5Str := hex.EncodeToString(cipherStr)
	return md5Str
}
