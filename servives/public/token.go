package public

import (
	"GoGameServer/core/libs/common"
	"GoGameServer/core/libs/dict"
	"GoGameServer/core/libs/jwt"
)

const (
	sessionSignKey            = "qwert&mnbvc"
	checkTimeOpen             = false
	userOffineCheckTime int64 = 10 * 60 * 1000
)

var (
	myJwt *jwt.Jwt
)

func init() {
	myJwt = jwt.NewJwt(sessionSignKey)
}

func CreateToken(userId uint64) string {
	claims := make(map[string]interface{})
	claims["userId"] = userId
	claims["time"] = common.UnixMillisecond()
	token := myJwt.Sign(claims)
	return token
}

func GetUserIdByToken(token string) uint64 {
	claims := myJwt.Parse(token)
	if claims == nil {
		return 0
	}

	//用户Id
	userId := dict.GetUint64(claims, "userId")

	//检测是否过期
	if checkTimeOpen {
		time := dict.GetInt64(claims, "time")
		nowTime := common.UnixMillisecond()
		if nowTime-time >= userOffineCheckTime {
			return 0
		}
	}

	return userId
}
