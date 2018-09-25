package public

import (
	"core/libs/common"
	"encoding/base64"
	"strconv"
	"strings"
)

var (
	sessionSignKey      string = "qwert&mnbvc"
	userOffineCheckTime int64  = 10 * 60 * 1000
)

func CreateToken(userId uint64) string {
	str := common.NumToString(userId) + "," + common.NumToString(common.UnixMillisecond()) + "," + sessionSignKey
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func GetUserIdByToken(token string) uint64 {
	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0
	}

	arr := strings.Split(string(decodeBytes), ",")
	if len(arr) != 3 {
		return 0
	}

	userId, _ := strconv.ParseUint(arr[0], 10, 64)
	//time, _ := strconv.ParseInt(arr[1], 10, 64)
	signKey := arr[2]

	//检测key值是否符合
	if signKey != sessionSignKey {
		return 0
	}

	////检测是否过期
	//nowTime := common.UnixMillisecond()
	//if nowTime-time >= userOffineCheckTime {
	//	return ""
	//}

	return userId
}
