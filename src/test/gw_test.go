package test

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

import (
	proto "code.google.com/p/goprotobuf/proto"
	"protos/gameProto"
	. "tools"
	"tools/cfg"
	"tools/unitest"
)

var protocol = packet.New(
	binary.SplitByUint32BE, 1024, 1024, 1024,
)

func Test_gateway(t *testing.T) {
	DEBUG("消息通信测试")
	var wg sync.WaitGroup
	var successNum uint32
	for i := 0; i < 3000; i++ {
		wg.Add(1)
		go func(flag int) {
			//			if flag != 0 && RandomInt31n(100) < 50 {
			//				flag -= 1
			//			}
			defer wg.Done()

			var count uint32 = 0
			var userName string = "User" + strconv.Itoa(flag)

			//超时处理
			timerFunc := func() {
				ERR("失败：", userName, count)
			}
			var timer *time.Timer = time.AfterFunc(10*time.Second, timerFunc)

			//连接服务器
			client, err := link.ConnectTimeout("tcp", "0.0.0.0:"+cfg.GetValue("gateway_port"), time.Second*3, protocol)
			if !unitest.NotError(t, err) {
				return
			}
			defer client.Close()

			count += 1

			//接收服务器连接成功消息
			var revMsg packet.RAW
			//			err = client.Receive(&revMsg)
			//			if !unitest.NotError(t, err) {
			//				return
			//			}
			//			msg := &gameProto.ConnectSuccessS2C{}
			//			proto.Unmarshal(revMsg[2:], msg)
			//			DEBUG(binary.GetUint16LE(revMsg[:2]), msg)

			count += 1

			//发送登录消息
			send_msg := createLoginBytes(userName)
			//			DEBUG("发送数据：", flag, send_msg)
			err = client.Send(send_msg)
			if !unitest.NotError(t, err) {
				return
			}

			count += 1

			//接受登录成功消息
			err = client.Receive(&revMsg)
			if !unitest.NotError(t, err) {
				return
			}
			//			DEBUG(binary.GetUint16LE(revMsg[:2]))
			msg1 := &gameProto.UserLoginS2C{}
			proto.Unmarshal(revMsg[2:], msg1)

			count += 1

			if !unitest.Pass(t, msg1.GetUserID() > 0) {
				return
			}

			count += 1

			//发送获取用户信息消息
			if msg1.GetUserID() != 0 {
				err = client.Send(createGetUserInfoBytes(msg1.GetUserID()))
				if !unitest.NotError(t, err) {
					return
				}

				count += 1

				//接受用户信息消息
				err = client.Receive(&revMsg)
				if !unitest.NotError(t, err) {
					return
				}

				count += 1

				if binary.GetUint16LE(revMsg[:2]) == gameProto.ID_ErrorMsgS2C {
					msg2 := &gameProto.ErrorMsgS2C{}
					proto.Unmarshal(revMsg[2:], msg2)
					//					DEBUG(binary.GetUint16LE(revMsg[:2]), msg2)
				} else {
					msg2 := &gameProto.GetUserInfoS2C{}
					proto.Unmarshal(revMsg[2:], msg2)
					//					DEBUG(binary.GetUint16LE(revMsg[:2]), msg2)

					successNum += 1
					DEBUG("成功：", userName, msg1.GetUserID(), successNum)
				}
			}

			timer.Stop()

		}(i)
	}
	wg.Wait()
}

func createLoginBytes(userName string) packet.RAW {
	sendMsg := &gameProto.UserLoginC2S{
		UserName: proto.String(userName),
	}
	msgBody, _ := proto.Marshal(sendMsg)

	msg1 := make([]byte, 2+len(msgBody))
	binary.PutUint16LE(msg1, gameProto.ID_UserLoginC2S)
	copy(msg1[2:], msgBody)
	return packet.RAW(msg1)
}

func createGetUserInfoBytes(userId uint64) packet.RAW {
	sendMsg := &gameProto.GetUserInfoC2S{
		UserID: proto.Uint64(userId),
	}
	msgBody, _ := proto.Marshal(sendMsg)

	msg1 := make([]byte, 2+len(msgBody))
	binary.PutUint16LE(msg1, gameProto.ID_GetUserInfoC2S)
	copy(msg1[2:], msgBody)
	return packet.RAW(msg1)
}

func RandBytes(n int) []byte {
	n = rand.Intn(n)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}
