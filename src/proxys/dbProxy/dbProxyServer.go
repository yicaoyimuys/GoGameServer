package dbProxy

import (
	"github.com/funny/link"
	"protos"
	"protos/systemProto"
	"protos/dbProto"
	"strings"
	. "tools"
	"tools/db"
	"tools/timer"
	"tools/dispatch"
	"proxys/redisProxy"
	"tools/debug"
	"global"
)

const (
	//处理写数据入库间隔(5分钟)
	SYSDB_INTERVAL = 5 * 60
)

var (
	servers                  map[string]*link.Session
	syncDbTimerID            uint64
	serverMsgReceiveChans 	 []dispatch.ReceiveMsgChan
	serverMsgDispatch        dispatch.DispatchInterface
	serverMsgDispatchAsync   dispatch.DispatchInterface
)

func init() {
	//创建异步接收消息的Chans
	serverMsgReceiveChans = make([]dispatch.ReceiveMsgChan, 10)
	for i := 0; i < len(serverMsgReceiveChans); i++ {
		serverMsgReceiveChans[i] = make(chan dispatch.ReceiveMsg, 2048)
	}

	//创建DB同步数据消息分派
	serverMsgDispatch = dispatch.NewDispatchAsync(serverMsgReceiveChans,
		dispatch.HandleCondition{
			Condition: dbProto.IsValidSyncID,
			H: dispatch.Handle{
				systemProto.ID_System_ConnectDBServerC2S:	connectDBServer,
				dbProto.ID_DB_User_LoginC2S:				userLogin,
			},
		},
	)

	//创建DB异步数据消息分派
	serverMsgDispatchAsync = dispatch.NewDispatch(
		dispatch.HandleCondition{
			Condition: dbProto.IsValidAsyncID,
			H: dispatch.Handle{
				dbProto.ID_DB_User_UpdateLastLoginTimeC2S:		updateUserLastLoginTime,
			},
		},
	)
}

//初始化
func InitServer(port string) error {
	servers = make(map[string]*link.Session)

	//开启DB
	db.Init()

	//开启同步写入DB
	StartSysDB()

	//监听tcp连接
	addr := "0.0.0.0:" + port
	err := global.Listener("tcp", addr, global.PackCodecType_Safe,
		func(session *link.Session) {},
		serverMsgDispatch,
	)

	return err
}

//客户端连接DBServer使用
func connectDBServer(session *link.Session, protoMsg protos.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectDBServerC2S)

	serverName := rev_msg.GetServerName()
	serverName = strings.Split(serverName, "[")[0]
	servers[serverName] = session

	session.AddCloseCallback(session, func(){
		delete(servers, serverName)
		ERR(serverName + " Disconnect At " + global.ServerName)
	})

	send_msg := protos.MarshalProtoMsg(&systemProto.System_ConnectDBServerS2C{})
	protos.Send(session, send_msg)
}

//开启定时同步DB数据
func StartSysDB() {
	syncDbTimerID = timer.DoTimer(int64(SYSDB_INTERVAL), onSyncDBTimer)
	onSyncDBTimer()
}

//停止定时同步DB数据
func SyncDB() {
	timer.Remove(syncDbTimerID)
	onSyncDBTimer()
}

//同步数据到DB服务器
func onSyncDBTimer() {
	debug.Start("SyncDBTimer")
	defer debug.Stop("SyncDBTimer")

	datas := redisProxy.PullDBWriteMsg()
	if datas == nil{
		return
	}
	dlen := len(datas)
	INFO("SyncDB Num: ", dlen)
	for i := 0; i < dlen; i++ {
		serverMsgDispatchAsync.Process(nil, datas[i])
	}
}

//发送DB消息到客户端
func sendDBMsgToClient(session *link.Session, msg []byte) {
	if session == nil {
		clientMsgDispatch.Process(session, msg)
	} else {
		protos.Send(session, msg)
	}
}
