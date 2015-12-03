package dbProxy

import (
	"container/list"
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"protos"
	"protos/dbProto"
	"protos/systemProto"
	"strings"
	. "tools"
	"tools/db"
	"tools/timer"
)

type goroutineMsg struct {
	msg     packet.RAW
	session *link.Session
}

type goroutineObj struct {
	revMsgChan chan goroutineMsg
}

const (
	SYSDB_INTERVAL = 5 * 60
)

var (
	servers                  map[string]*link.Session
	asyncMsgs                *list.List
	syncDbTimerID            uint64
	revSyncMsgGoroutines     []goroutineObj
	revSyncMsgGoroutineIndex int
	revSyncMsgGoroutineNum   int
)

//初始化
func InitServer(port string) error {
	servers = make(map[string]*link.Session)
	asyncMsgs = list.New()

	db.Init()

	startSysDB()

	createRevGoroutines()

	listener, err := link.Serve("tcp", "0.0.0.0:"+port, packet.New(
		binary.SplitByUint32BE, 1024, 1024, 1024,
	))
	if err != nil {
		return err
	}

	listener.Serve(func(session *link.Session) {
		for {
			var msg packet.RAW
			if err := session.Receive(&msg); err != nil {
				break
			}

			msgID := binary.GetUint16LE(msg[:2])
			if systemProto.IsValidID(msgID) {
				//系统消息
				dealReceiveSystemMsgC2S(session, msg)
			} else if dbProto.IsValidAsyncID(msgID) {
				//异步DB消息
				asyncMsgs.PushBack(msg)
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

	return nil
}

//处理接收到的系统消息
func dealReceiveSystemMsgC2S(session *link.Session, msg packet.RAW) {
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

	send_msg := systemProto.MarshalProtoMsg(&systemProto.System_ConnectDBServerS2C{})
	protos.Send(send_msg, session)
}

//开启定时同步DB数据
func startSysDB() {
	syncDbTimerID = timer.DoTimer(int64(SYSDB_INTERVAL), onSyncDBTimer)
}

//停止定时同步DB数据
func StopSyncDB() {
	timer.Remove(syncDbTimerID)
	onSyncDBTimer()
}

//同步数据到DB服务器
func onSyncDBTimer() {
	INFO("SyncDB Num: ", asyncMsgs.Len())
	for msg := asyncMsgs.Front(); msg != nil; msg = msg.Next() {
		protoMsg := msg.Value.(packet.RAW)
		dealReceiveAsyncDBMsgC2S(protoMsg)
	}
	asyncMsgs.Init()
}

//创建接收同步消息的Goroutines
func createRevGoroutines() {
	revSyncMsgGoroutineIndex = 0
	revSyncMsgGoroutineNum = 5
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
