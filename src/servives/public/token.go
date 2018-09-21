package public

import (
	"core/libs/common"
	"encoding/base64"
	"strings"
)

var (
	sessionSignKey      string = "qwert&mnbvc"
	userOffineCheckTime int64  = 10 * 60 * 1000
)

func CreateToken(userId string) string {
	str := userId + "," + common.NumToString(common.UnixMillisecond()) + "," + sessionSignKey
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func GetUserIdByToken(token string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return ""
	}

	arr := strings.Split(string(decodeBytes), ",")
	if len(arr) != 3 {
		return ""
	}

	userId := arr[0]
	//time, _ := strconv.ParseInt(arr[1], 10, 64)
	signKey := arr[2]

	//检测key值是否符合
	if signKey != sessionSignKey {
		return ""
	}

	////检测是否过期
	//nowTime := common.UnixMillisecond()
	//if nowTime-time >= userOffineCheckTime {
	//	return ""
	//}

	return userId
}
