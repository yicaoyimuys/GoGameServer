package main

import (
	"core/consts/errCode"
	"core/consts/service"
	. "core/libs"
	"core/libs/array"
	"core/libs/hash"
	"core/libs/random"
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
	"sync"
	"sync/atomic"
)

var (
	servers = []string{
		"127.0.0.1:19881",
		"127.0.0.1:19882",
	}

	connectNums  = 100
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
	for {
		account := "ys" + NumToString(random.RandIntn(10000))
		if array.InArray(userAccounts, account) {
			continue
		}
		userAccounts = append(userAccounts, account)
		if len(userAccounts) == connectNums {
			break
		}
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
	client.login()
}

type clientSession struct {
	con         *websocket.Conn
	account     string
	token       string
	pingTimerId *timer.TimerEvent
	closeFlag   int32

	recvMutex sync.Mutex
	sendMutex sync.Mutex
}

//关闭
func (this *clientSession) close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		timer.Remove(this.pingTimerId)
		this.con.Close()
	}
}

//是否关闭
func (this *clientSession) isClose() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

//平台登录
func (this *clientSession) login() {
	msg := &gameProto.UserLoginC2S{
		Account: protos.String(this.account),
	}
	this.sendMsg(msg)
}

//获取用户数据
func (this *clientSession) getInfo() {
	msg := &gameProto.UserGetInfoC2S{
		Token: protos.String(this.token),
	}
	this.sendMsg(msg)
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

	if msgId == gameProto.ID_user_login_s2c {
		//登录成功
		data := msgData.(*gameProto.UserLoginS2C)
		this.token = data.GetToken()
		DEBUG("登录成功", this.token)
		this.getInfo()
	} else if msgId == gameProto.ID_user_getInfo_s2c {
		//获取用户信息成功
		data := msgData.(*gameProto.UserGetInfoS2C)
		DEBUG("用户信息", data.GetData())
	} else if msgId == gameProto.ID_error_notice_s2c {
		data := msgData.(*gameProto.ErrorNoticeS2C)
		errCode := data.GetErrorCode()
		if errCode == ErrCode.SYSTEM_ERR {
			//系统服务错误
			this.close()
			//重新连接
			timer.SetTimeOut(3000, func() {
				go startConnect(this.account)
			})
		}
		DEBUG("收到错误消息", data)
	}
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
