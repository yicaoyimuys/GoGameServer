package global

import (
//	"runtime"
	. "tools"
	. "tools/gc"
	"tools/codecType"
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
	SysRoutine()

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
func Listener(network, address string, codecType link.CodecType, acceptFunc func(*link.Session)) error {
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
		}
	}()
	return nil
}
