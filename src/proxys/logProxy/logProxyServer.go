package logProxy

import (
	"github.com/funny/link"
	"protos"
	"protos/logProto"
	"protos/systemProto"
	"time"
	"strconv"
	. "tools"
	"tools/file"
	"global"
	"tools/dispatch"
)

var (
	servers map[string]*link.Session
	serverMsgReceiveChan dispatch.ReceiveMsgChan
	serverMsgDispatch dispatch.DispatchInterface
)

func init() {
	//创建异步接收消息的Chan
	serverMsgReceiveChan = make(chan dispatch.ReceiveMsg, 4096)

	//创建消息分派
	serverMsgDispatch = dispatch.NewDispatchAsync([]dispatch.ReceiveMsgChan{serverMsgReceiveChan},
		dispatch.Handle{
			systemProto.ID_System_ConnectLogServerC2S:	connectLogServer,
			logProto.ID_Log_CommonLogC2S:				writeLogFile,
		},
	)
}

//初始化
func InitServer(port string) error {
	servers = make(map[string]*link.Session)

	//监听tcp连接
	addr := "0.0.0.0:" + port
	err := global.Listener("tcp", addr, global.PackCodecType_Safe,
		func(session *link.Session) {},
		serverMsgDispatch,
	)

	return err
}

//等待所有log写入文件
func SyncLog() {
	INFO("SyncLog Num: ", len(serverMsgReceiveChan))
	for len(serverMsgReceiveChan) > 0 {
		
	}
	close(serverMsgReceiveChan)
}

//写入log
func writeLogFile(session *link.Session, msg protos.ProtoMsg)  {
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
func connectLogServer(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectLogServerC2S)

	serverName := rev_msg.GetServerName()
	servers[serverName] = session

	session.AddCloseCallback(session, func(){
		delete(servers, serverName)
		ERR(serverName + " Disconnect At " + global.ServerName)
	})

	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectLogServerS2C{})
	protos.Send(session, send_msg)
}
