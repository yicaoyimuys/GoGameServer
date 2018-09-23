package main

import (
	"core/consts/service"
	. "core/libs"
	"core/libs/hash"
	"core/libs/stack"
	"core/libs/timer"
	"core/protos"
	"core/protos/gameProto"
	"core/service"
	"crypto/tls"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	_ "net/http/pprof"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	servers = []string{
		"127.0.0.1:19881",
		"127.0.0.1:19882",
	}

	connectNums  = 10
	userAccounts = []string{}
)

func main() {
	//初始化Service
	service.NewService(Service.Test)

	//开始测试
	startTest()

	//保持进程
	Run()
}

func startTest() {
	//初始化使用账号
	initAccount()

	//开启链接
	for i := 0; i < connectNums; i++ {
		go startConnect(userAccounts[i])
	}
}

func initAccount() {
	for a := 0; a < connectNums; a++ {
		account := strconv.Itoa(10000 + a)
		userAccounts = append(userAccounts, account)
	}
}

func startConnect(account string) {
	hash := hash.GetHash([]byte(account))

	var serverIndex = hash % uint32(len(servers))
	var server = servers[serverIndex]
	//DEBUG(myUserId, serverIndex)

	scheme := "ws"
	u := url.URL{Scheme: scheme, Host: server, Path: "/"}
	d := websocket.DefaultDialer
	d.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := d.Dial(u.String(), nil)
	if err != nil {
		ERR(account, "连接失败")
		return
	}

	INFO(account, "连接成功")

	client := new(clientSession)
	client.con = c
	client.account = account

	go client.receiveMsg()

	client.ping()
	client.network()
	client.platformLogin()
}

type clientSession struct {
	con             *websocket.Conn
	account         string
	progressNum     uint16
	progressTimerId *timer.TimerEvent
	eventTimerId    *timer.TimerEvent
	pingTimerId     *timer.TimerEvent
	networkTimerId  *timer.TimerEvent
	resultTimerId   *timer.TimerEvent
	closeFlag       int32

	recvMutex sync.Mutex
	sendMutex sync.Mutex
}

//关闭
func (this *clientSession) close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		timer.Remove(this.progressTimerId)
		timer.Remove(this.eventTimerId)
		timer.Remove(this.pingTimerId)
		timer.Remove(this.networkTimerId)
		timer.Remove(this.resultTimerId)
		this.con.Close()
	}
}

//是否关闭
func (this *clientSession) isClose() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

//平台登录
func (this *clientSession) platformLogin() {
	//userId := this.myUserId
	//otherUserId := this.otherUserId
	//roomId := this.roomId
	//
	//platformData := make(map[string]interface{})
	//platformData["userId"] = userId
	//platformData["userName"] = "Name" + userId
	//platformData["userPic"] = "Pic" + userId
	//platformData["userSex"] = 1
	//
	//platformData["otherUserId"] = otherUserId
	//platformData["otherUserName"] = "Name" + otherUserId
	//platformData["otherUserPic"] = "Pic" + otherUserId
	//platformData["otherUserSex"] = 2
	//
	//platformData["roomId"] = roomId
	//platformData["isAi"] = 0
	//platformData["aiLevel"] = 0
	//
	//platformDataStr, _ := json.Marshal(platformData)
	//
	//var sendMsg = msg.NewPlatform_login_c2s()
	//sendMsg.GameId = 1
	//sendMsg.PlatformId = 0
	//sendMsg.PlatformData = string(platformDataStr)
	//this.sendMsg(sendMsg)
}

//进入游戏
func (this *clientSession) joinGame() {
	//var sendMsg = msg.NewGame_join_c2s()
	//sendMsg.RoomId = this.roomId
	//this.sendMsg(sendMsg)
}

//发送进度
func (this *clientSession) loadProgress() {
	//this.progressTimerId = timer.DoTimer(1000, func() {
	//	this.progressNum += uint16(random.RandIntRange(1, 50))
	//	if this.progressNum >= 100 {
	//		this.progressNum = 100
	//		timer.Remove(this.progressTimerId)
	//	}
	//
	//	var msg = msg.NewGame_loadProgress_c2s()
	//	msg.RoomId = this.roomId
	//	msg.Progress = this.progressNum
	//	this.sendMsg(msg)
	//})
}

//发送游戏事件
func (this *clientSession) gameEvent() {
	//this.eventTimerId = timer.DoTimer(500, func() {
	//	msg := msg.NewGame_event_c2s()
	//	msg.RoomId = this.roomId
	//	msg.SendType = 1
	//	msg.Event = "哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈" + this.myUserId
	//	this.sendMsg(msg)
	//})
}

//心跳
func (this *clientSession) ping() {
	this.pingTimerId = timer.DoTimer(2000, func() {
		var msg = &gameProto.ClientPingC2S{}
		this.sendMsg(msg)
	})
}

//网络监测
func (this *clientSession) network() {
	//this.networkTimerId = timer.DoTimer(5000, func() {
	//	var msg = msg.NewClient_network_c2s()
	//	msg.Time = uint32(time.Now().Unix())
	//	this.sendMsg(msg)
	//})
}

//结算
func (this *clientSession) gameResult() {
	//msg := msg.NewGame_result_c2s()
	//msg.RoomId = this.roomId
	//msg.Result = 1
	//this.sendMsg(msg)
}

//接收消息
func (this *clientSession) receiveMsg() {
	defer stack.PrintPanicStackError()

	for {
		if this.isClose() {
			break
		}

		_, message, err := this.con.ReadMessage()
		if err != nil {
			ERR("read:", err)
			break
		}

		//消息内容
		msgBody := message[2:]
		//消息解析
		protoMsg := protos.UnmarshalProtoMsg(msgBody)
		if protoMsg == protos.NullProtoMsg {
			ERR("收到错误消息ID: ", protos.UnmarshalProtoId(msgBody))
			break
		}
		DEBUG(this.account, "收到消息ID", protoMsg.ID)
		//消息处理
		this.handleMsg(protoMsg.ID, protoMsg.Body)
	}
	this.close()
}

func (this *clientSession) handleMsg(msgId uint16, msgData proto.Message) {
	defer stack.PrintPanicStackError()

	//if msgId == msg.ID_Platform_login_s2c {
	//	this.joinGame()
	//} else if msgId == msg.ID_Game_join_s2c {
	//	this.loadProgress()
	//} else if msgId == msg.ID_Game_loadProgress_notice_s2c {
	//	data := msgData.(*msg.Game_loadProgress_notice_s2c)
	//	DEBUG(this.myUserId, "收到加载进度", *data)
	//} else if msgId == msg.ID_Game_start_notice_s2c {
	//	data := msgData.(*msg.Game_start_notice_s2c)
	//	DEBUG("收到游戏开启", *data)
	//	//发送游戏事件
	//	this.gameEvent()
	//	//发送结算
	//	arr := strings.Split(this.roomId, "_")
	//	if this.myUserId == arr[0] {
	//		this.resultTimerId = timer.SetTimeOut(60*1000, func() {
	//			this.gameResult()
	//		})
	//	}
	//} else if msgId == msg.ID_Game_event_notice_s2c {
	//	data := msgData.(*msg.Game_event_notice_s2c)
	//	DEBUG("收到游戏事件", *data)
	//} else if msgId == msg.ID_Client_network_s2c {
	//	data := msgData.(*msg.Client_network_s2c)
	//	cha := uint32(time.Now().Unix()) - data.Time
	//	if cha > 100 {
	//		DEBUG("延迟：", cha)
	//	}
	//} else if msgId == msg.ID_Game_result_notice_s2c {
	//	data := msgData.(*msg.Game_result_notice_s2c)
	//	DEBUG("收到游戏结束", *data)
	//	this.close()
	//
	//	//重新连接
	//	timer.SetTimeOut(2000, func() {
	//		startConnect(this.myUserId, this.otherUserId)
	//	})
	//} else if msgId == msg.ID_Error_notice_s2c {
	//	data := msgData.(*msg.Error_notice_s2c)
	//	DEBUG("收到错误消息", *data)
	//}
}

func (this *clientSession) sendMsg(msg proto.Message) {
	if this.isClose() {
		return
	}

	defer stack.PrintPanicStackError()

	this.sendMutex.Lock()
	defer this.sendMutex.Unlock()

	msgBytes := protos.MarshalProtoMsg(msg)

	msgLen := uint16(len(msgBytes))
	sendMsg := make([]byte, msgLen+2)
	binary.BigEndian.PutUint16(sendMsg[:2], msgLen)
	copy(sendMsg[2:], msgBytes)

	this.con.WriteMessage(websocket.BinaryMessage, sendMsg)
}
