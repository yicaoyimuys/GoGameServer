package test

import (
	"github.com/funny/binary"
	"github.com/funny/link"
	"github.com/funny/link/packet"
	"github.com/funny/unitest"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

import (
	proto "code.google.com/p/goprotobuf/proto"
	"protos"
	. "tools"
	"tools/cfg"
)

var protocol = packet.New(
	binary.SplitByUint32BE, 1024, 1024, 1024,
)

func Test_gateway(t *testing.T) {
	DEBUG("消息通信测试")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(flag int) {
			defer wg.Done()

			client, err := link.Connect("tcp", "120.25.202.132:"+cfg.GetValue("gateway_port"), protocol)
			if !unitest.NotError(t, err) {
				return
			}

			defer client.Close()

			//接受服务器连接成功消息
			var revMsg packet.RAW
			err = client.Receive(&revMsg)
			if !unitest.NotError(t, err) {
				return
			}

			msg := &simple.ConnectSuccessS2C{}
			proto.Unmarshal(revMsg[2:], msg)

			DEBUG(binary.GetUint16LE(revMsg[:2]), msg)

			//发送登录消息
			err = client.Send(createLoginBytes("User" + strconv.Itoa(flag)))
			if !unitest.NotError(t, err) {
				return
			}

			//接受登录成功消息
			err = client.Receive(&revMsg)
			if !unitest.NotError(t, err) {
				return
			}

			msg1 := &simple.UserLoginS2C{}
			proto.Unmarshal(revMsg[2:], msg1)

			DEBUG(binary.GetUint16LE(revMsg[:2]), msg1)

			//发送获取用户信息消息
			if msg1.GetUserID() != -1 {
				err = client.Send(createGetUserInfoBytes(msg1.GetUserID()))
				if !unitest.NotError(t, err) {
					return
				}

				//接受用户信息消息
				err = client.Receive(&revMsg)
				if !unitest.NotError(t, err) {
					return
				}

				if binary.GetUint16LE(revMsg[:2]) == simple.ID_ErrorMsg {
					msg2 := &simple.ErrorMsg{}
					proto.Unmarshal(revMsg[2:], msg2)

					DEBUG(binary.GetUint16LE(revMsg[:2]), msg2)
				} else {
					msg2 := &simple.GetUserInfoS2C{}
					proto.Unmarshal(revMsg[2:], msg2)

					DEBUG(binary.GetUint16LE(revMsg[:2]), msg2)
				}
			}

		}(i)
	}
	wg.Wait()

}

func createLoginBytes(userName string) packet.RAW {
	sendMsg := &simple.UserLoginC2S{
		UserName: proto.String(userName),
	}
	msgBody, _ := proto.Marshal(sendMsg)

	msg1 := make([]byte, 2+len(msgBody))
	binary.PutUint16LE(msg1, simple.ID_UserLoginC2S)
	copy(msg1[2:], msgBody)
	return packet.RAW(msg1)
}

func createGetUserInfoBytes(userId int32) packet.RAW {
	sendMsg := &simple.GetUserInfoC2S{
		UserID: proto.Int32(userId),
	}
	msgBody, _ := proto.Marshal(sendMsg)

	msg1 := make([]byte, 2+len(msgBody))
	binary.PutUint16LE(msg1, simple.ID_GetUserInfoC2S)
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
