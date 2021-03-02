package timer

import (
	"GoGameServer/core/libs/stack"
	"sync/atomic"
	"time"
)

type TimerEvent struct {
	callBack    func()       // 回调函数
	delay       uint32       // 执行间隔
	repeatCount uint32       // 执行次数
	ticker      *time.Ticker //定时器
	closeFlag   int32        //关闭标识
	closeChan   chan int     //关闭使用
}

//是否已关闭
func (this *TimerEvent) IsClosed() bool {
	return atomic.LoadInt32(&this.closeFlag) == 1
}

//关闭
func (this *TimerEvent) Close() {
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		if this.ticker != nil {
			this.ticker.Stop()
		}
		close(this.closeChan)
	}
}

//无限次数执行
func DoTimer(delay uint32, callback func()) *TimerEvent {
	return Do(delay, 0, callback)
}

//延时处理
func SetTimeOut(delay uint32, callback func()) *TimerEvent {
	return Do(delay, 1, callback)
}

//移除一个定时器
func Remove(event *TimerEvent) {
	if event == nil {
		return
	}
	event.Close()
}

//时间间隔，执行次数，回调函数
func Do(delay uint32, repeatCount uint32, callback func()) *TimerEvent {
	//最小单位1ms
	if delay < 1 {
		callback()
		return nil
	}

	//创建事件对象
	event := &TimerEvent{
		callBack:    callback,
		delay:       delay,
		repeatCount: repeatCount,
		closeChan:   make(chan int),
	}

	//开启timer
	go startTicker(event)

	//返回
	return event
}

func startTicker(event *TimerEvent) {
	defer stack.TryError()
	event.ticker = time.NewTicker(time.Duration(event.delay) * time.Millisecond)
	for {
		select {
		case <-event.ticker.C:
			event.callBack()
			if event.repeatCount > 0 {
				event.repeatCount -= 1
				if event.repeatCount == 0 {
					Remove(event)
				}
			}
		case <-event.closeChan:
			return
		}
	}
}
