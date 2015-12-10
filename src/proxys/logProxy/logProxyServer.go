package logProxy

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"protos"
	"protos/logProto"
	"protos/systemProto"
	"time"
	"strconv"
	. "tools"
	"tools/file"
)

type receiveMsg struct {
	session *link.Session
	msg     packet.RAW
}

var (
	servers map[string]*link.Session
	receiveMsgs chan receiveMsg
)

//初始化
func InitServer(port string) error {
	servers = make(map[string]*link.Session)
	receiveMsgs = make(chan receiveMsg, 2048)

	listener, err := link.Serve("tcp", "0.0.0.0:" + port, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	if err != nil {
		return err
	}

	go func() {
		listener.Serve(func(session *link.Session) {
			for {
				var msg packet.RAW
				if err := session.Receive(&msg); err != nil {
					break
				}
				receiveMsgs <- receiveMsg{
					session:session,
					msg:msg,
				}
			}
		})
	}()

	go dealReceiveMsgs()

	return nil
}

func dealReceiveMsgs() {
	for {
		data, ok := <-receiveMsgs
		if !ok {
			return
		}
		dealReceiveMsgC2S(data.session, data.msg)
	}
}

//处理接收到的消息
func dealReceiveMsgC2S(session *link.Session, msg packet.RAW) {
	if len(msg) < 2 {
		return
	}

	msgID := binary.GetUint16LE(msg[:2])
	if systemProto.IsValidID(msgID) {
		//系统消息
		dealReceiveSystemMsgC2S(session, msg)
	} else if logProto.IsValidID(msgID) {
		//Log消息
		dealLogMsgC2S(msg)
	}
}

//处理接收到的系统消息
func dealReceiveSystemMsgC2S(session *link.Session, msg packet.RAW) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectLogServerC2S:
		connectLogServer(session, protoMsg)
	}
}

//处理Log消息
func dealLogMsgC2S(msg packet.RAW) {
	protoMsg := logProto.UnmarshalProtoMsg(msg)
	if protoMsg == logProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case logProto.ID_Log_CommonLogC2S:
		writeLogFile(protoMsg)
	}
}

//等待所有log写入文件
func SyncLog() {
	INFO("SyncLog Num: ", len(receiveMsgs))
	for len(receiveMsgs) > 0 {
		
	}
}

//写入log
func writeLogFile(msg logProto.ProtoMsg)  {
	data := msg.Body.(*logProto.Log_CommonLogC2S)

	t := time.Unix(data.GetTime(), 0)

	//创建目录
	dirPath := "gamelogs/" + data.GetDir() + "/"+ t.Format("2006-01-02")
	err := file.CreateDir(dirPath)
	if err != nil {
		return
	}

	//创建文件
	filePath := dirPath + "/" + t.Format("15") + ".log";
	file := file.OpenFile(filePath)
	if file == nil{
		return
	}
	
	defer file.Close()

	//写入文件
	str := ""
	str += strconv.FormatUint(uint64(data.GetType()), 10)
	str += "_"
	str += strconv.FormatInt(data.GetTime(), 10)
	str += "_"
	str += data.GetContent()
	str += "\n"
	file.WriteString(str)
}

//其他客户端连接LogServer处理
func connectLogServer(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectLogServerC2S)

	serverName := rev_msg.GetServerName()
	servers[serverName] = session

	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectLogServerS2C{})
	protos.Send(send_msg, session)
}
