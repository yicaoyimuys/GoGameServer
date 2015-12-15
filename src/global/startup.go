package global

import (
//	"runtime"
	. "tools"
	"tools/gc"
	"tools/codecType"
	"tools/dispatch"
	"github.com/funny/link"
)

var (
	PackCodecType_UnSafe 	link.CodecType = link.Packet(4, 1024 * 1024, 4096, link.LittleEndian, codecType.NetCodecType{})
	PackCodecType_Safe 		link.CodecType = link.ThreadSafe(PackCodecType_UnSafe)
	PackCodecType_Async 	link.CodecType = link.Async(4096, PackCodecType_UnSafe)
)

//服务器启动
func Startup(serverName string, logFile string, stopServerFunc func()) {
	//	runtime.GOMAXPROCS(runtime.NumCPU())

	// 开启Log记录
	SetLogFile(logFile)
	SetLogPrefix(serverName)

	// 信号量监听
	go SignalProc(stopServerFunc)

	// 开启GC及系统环境信息监测
	gc.SysRoutine()

	// 开启服务器
	INFO("Starting " + serverName)
}

// 保持进程
func Run() {
	temp := make(chan int32, 10)
	for {
		select {
		case <-temp:
		}
	}
}

// 开启服务器监听
func Listener(network, address string, codecType link.CodecType, acceptFunc func(*link.Session), dispatch dispatch.DispatchInterface) error {
	listener, err := link.Serve(network, address, codecType)
	if err != nil {
		return err
	}

	go func() {
		for {
			session, err := listener.Accept()
			if err != nil {
				break
			}

			go acceptFunc(session)

			if dispatch != nil{
				go sessionReceive(session, dispatch)
			}
		}
	}()
	return nil
}

// 连接服务器
func Connect(connectServerName string, network string, address string, codecType link.CodecType, dispatch dispatch.DispatchInterface) (*link.Session, error) {
	session, err := link.Connect(network, address, codecType)
	if err == nil {
		session.AddCloseCallback(session, func(){
			ERR(connectServerName + " Disconnect At " + ServerName)
		})

		if dispatch != nil{
			go sessionReceive(session, dispatch)
		}
	}
	return session, err
}

func sessionReceive(session *link.Session, d dispatch.DispatchInterface) {
	for {
		var msg []byte
		if err := session.Receive(&msg); err != nil {
			break
		}

		d.Process(session, msg)
	}
}
