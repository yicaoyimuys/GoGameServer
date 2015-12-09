package logProxy

import (
	"protos/logProto"
	"protos"
	"time"
	"strconv"
)

const (
	Type_UserLogin = 1
)

func sendCommonLog(dir string, logType uint32, content string)  {
	send_msg := logProto.MarshalProtoMsg(&logProto.Log_CommonLogC2S{
		Dir: protos.String(dir),
		Type: protos.Uint32(logType),
		Content: protos.String(content),
		Time: protos.Int64(time.Now().Unix()),
	})
	SendLogMsgToServer(send_msg)
}

func sendLoginLog(logType uint32, content string)  {
	sendCommonLog("login", logType, content)
}

func UserLogin(userID uint64)  {
	sendLoginLog(Type_UserLogin, strconv.FormatUint(userID, 10))
}