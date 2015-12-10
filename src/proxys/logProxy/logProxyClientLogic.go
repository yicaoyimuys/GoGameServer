package logProxy

import (
	"protos/logProto"
	"protos"
	"time"
	"strconv"
	"strings"
)

const (
	type_UserLogin = 1
	type_UserOffLine = 2
)

func sendCommonLog(dir string, logType uint32, content string) {
	send_msg := logProto.MarshalProtoMsg(&logProto.Log_CommonLogC2S{
		Dir: protos.String(dir),
		Type: protos.Uint32(logType),
		Content: protos.String(content),
		Time: protos.Int64(time.Now().Unix()),
	})
	SendLogMsgToServer(send_msg)
}

func sendLoginLog(logType uint32, contents []string) {
	sendCommonLog("login", logType, strings.Join(contents, "_"))
}

func UserLogin(userID uint64) {
	contents := []string{strconv.FormatUint(userID, 10)}
	sendLoginLog(type_UserLogin, contents)
}

func UserOffLine(userID uint64)  {
	contents := []string{strconv.FormatUint(userID, 10)}
	sendLoginLog(type_UserOffLine, contents)
}