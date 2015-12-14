package dbProxy

import (
	"github.com/funny/link"
	"github.com/funny/binary"
	"protos"
	"protos/systemProto"
	"strings"
	. "tools"
	"tools/db"
	"tools/timer"
	"proxys/redisProxy"
	"tools/debug"
	"global"
)

type goroutineMsg struct {
	msg     []byte
	session *link.Session
}

type goroutineObj struct {
	revMsgChan chan goroutineMsg
}

const (
	//处理写数据入库间隔(5分钟)
	SYSDB_INTERVAL = 5 * 60
)

var (
	servers                  map[string]*link.Session
	syncDbTimerID            uint64
	revSyncMsgGoroutines     []goroutineObj
	revSyncMsgGoroutineIndex int
	revSyncMsgGoroutineNum   int = 10
)

//初始化
func InitServer(port string) error {
	servers = make(map[string]*link.Session)

	db.Init()

	startSysDB()

	createRevGoroutines()

	err := global.Listener("tcp", "0.0.0.0:"+port, global.PackCodecType, func(session *link.Session) {
		for {
			var msg []byte
			if err := session.Receive(&msg); err != nil {
				break
			}

			msgID := binary.GetUint16LE(msg[:2])
			if systemProto.IsValidID(msgID) {
				//系统消息
				dealReceiveSystemMsgC2S(session, msg)
			} else {
				//同步DB消息
				useObj := revSyncMsgGoroutines[revSyncMsgGoroutineIndex]
				useObj.revMsgChan <- goroutineMsg{msg, session}
				revSyncMsgGoroutineIndex++
				if revSyncMsgGoroutineIndex == revSyncMsgGoroutineNum {
					revSyncMsgGoroutineIndex = 0
				}
			}
		}
	})

	if err != nil {
		return err
	}

	return nil
}

//处理接收到的系统消息
func dealReceiveSystemMsgC2S(session *link.Session, msg []byte) {
	protoMsg := systemProto.UnmarshalProtoMsg(msg)
	if protoMsg == systemProto.NullProtoMsg {
		return
	}

	switch protoMsg.ID {
	case systemProto.ID_System_ConnectDBServerC2S:
		connectDBServer(session, protoMsg)
	}
}

//客户端连接DBServer使用
func connectDBServer(session *link.Session, protoMsg systemProto.ProtoMsg) {
	rev_msg := protoMsg.Body.(*systemProto.System_ConnectDBServerC2S)

	serverName := rev_msg.GetServerName()
	serverName = strings.Split(serverName, "[")[0]
	servers[serverName] = session

	session.AddCloseCallback(session, func(){
		delete(servers, serverName)
		ERR(serverName + " Disconnect At " + global.ServerName)
	})

	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectDBServerS2C{})
	protos.Send(send_msg, session)
}

//开启定时同步DB数据
func startSysDB() {
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
		msg := datas[i]
		dealReceiveAsyncDBMsgC2S(msg)
	}
}

//创建接收同步消息的Goroutines
func createRevGoroutines() {
	revSyncMsgGoroutineIndex = 0
	revSyncMsgGoroutines = make([]goroutineObj, revSyncMsgGoroutineNum)

	for i := 0; i < revSyncMsgGoroutineNum; i++ {
		obj := goroutineObj{
			revMsgChan: make(chan goroutineMsg, 2048),
		}

		revSyncMsgGoroutines[i] = obj
		go func() {
			for {
				goroutineMsg, ok := <-obj.revMsgChan
				if !ok {
					return
				}
				dealReceiveDBMsgC2S(goroutineMsg.session, goroutineMsg.msg)
			}
		}()
	}
}

//发送DB消息到客户端
func sendDBMsgToClient(session *link.Session, msg []byte) {
	if session == nil {
		dealReceiveDBMsgS2C(msg)
	} else {
		protos.Send(msg, session)
	}
}
