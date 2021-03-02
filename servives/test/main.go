package main

import (
	"GoGameServer/core/consts/ErrCode"
	"GoGameServer/core/consts/Service"
	. "GoGameServer/core/libs"
	"GoGameServer/core/libs/array"
	"GoGameServer/core/libs/hash"
	"GoGameServer/core/libs/random"
	"GoGameServer/core/libs/stack"
	"GoGameServer/core/libs/timer"
	"GoGameServer/core/protos"
	"GoGameServer/core/protos/gameProto"
	"GoGameServer/core/service"
	"crypto/tls"
	"encoding/binary"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	servers = []string{
		"127.0.0.1:19881",
		"127.0.0.1:19882",
	}

	connectNum   = 100
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
	for i := 0; i < len(userAccounts); i++ {
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
		if len(userAccounts) == connectNum {
			break
		}
	}
}

func startConnect(account string) {
	hashCode := hash.GetHash([]byte(account))

	var serverIndex = hashCode % uint32(len(servers))
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
	client.login()
}

type clientSession struct {
	con         *websocket.Conn
	account     string
	token       string
	pingTimerId *timer.TimerEvent
	chatTimerId *timer.TimerEvent
	closeFlag   int32

	recvMutex sync.Mutex
	sendMutex sync.Mutex
}

//关闭
func (this *clientSession) close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		timer.Remove(this.pingTimerId)
		timer.Remove(this.chatTimerId)
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

//加入聊天
func (this *clientSession) joinChat() {
	var msg = &gameProto.UserJoinChatC2S{}
	msg.Token = protos.String(this.token)
	this.sendMsg(msg)
}

//聊天
func (this *clientSession) chat() {
	delay := random.RandIntRange(10000, 20000)
	this.chatTimerId = timer.DoTimer(uint32(delay), func() {
		var msg = &gameProto.UserChatC2S{}
		msg.Msg = protos.String("你好啊")
		this.sendMsg(msg)
	})
}

//接收消息
func (this *clientSession) receiveMsg() {
	defer stack.TryError()

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
	defer stack.TryError()

	if msgId == gameProto.ID_user_login_s2c {
		//登录成功
		data := msgData.(*gameProto.UserLoginS2C)
		this.token = data.GetToken()
		DEBUG("登录成功", this.token)
		//获取用户数据
		this.getInfo()
	} else if msgId == gameProto.ID_user_getInfo_s2c {
		//获取用户信息成功
		data := msgData.(*gameProto.UserGetInfoS2C)
		DEBUG("用户信息", data.GetData())
		//加入聊天
		this.joinChat()
	} else if msgId == gameProto.ID_user_joinChat_s2c {
		//加入聊天成功
		DEBUG("加入聊天成功")
		//开始聊天
		this.chat()
	} else if msgId == gameProto.ID_user_chat_notice_s2c {
		//收到聊天消息
		data := msgData.(*gameProto.UserChatNoticeS2C)
		DEBUG(this.account, "收到聊天消息", data.GetUserName()+"说："+data.GetMsg())
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

	defer stack.TryError()

	this.sendMutex.Lock()
	defer this.sendMutex.Unlock()

	msgBytes := protos.MarshalProtoMsg(msg)

	msgLen := uint16(len(msgBytes))
	sendMsg := make([]byte, msgLen+2)
	binary.BigEndian.PutUint16(sendMsg[:2], msgLen)
	copy(sendMsg[2:], msgBytes)

	this.con.WriteMessage(websocket.BinaryMessage, sendMsg)
}
