package timer

import (
	"container/list"
	"sync/atomic"
	"time"
	. "tools"
)

type timerEvent struct {
	ID          uint64 // TimerID
	Func        func() // 回调函数
	Delay       int64  // 执行间隔
	RepeatCount int    // 执行次数
}

var (
	eventQueues  map[int64]*list.List
	events       map[uint64]*timerEvent
	maxSessionId uint64
)

func init() {
	eventQueues = make(map[int64]*list.List)
	events = make(map[uint64]*timerEvent)
	maxSessionId = 0
	go timer()
}

func timer() {
	defer func() {
		if x := recover(); x != nil {
			ERR("TIMER CRASHED", x)
		}
	}()

	nowTime := time.Now().Unix()
	sleep_timer := time.NewTimer(time.Second)
	for {
		select {
		case <-sleep_timer.C:
			sleep_timer.Reset(time.Second)
			nowTime += 1
			queues, exists := eventQueues[nowTime]
			if exists {
				for evt := queues.Front(); evt != nil; evt = evt.Next() {
					eventID := evt.Value.(uint64)
					if event, existsEvent := events[eventID]; existsEvent {
						event.Func()
						if event.RepeatCount > 0 {
							event.RepeatCount -= 1
							if event.RepeatCount == 0 {
								delete(events, eventID)
							} else {
								addToQueue(nowTime, event)
							}
						} else {
							addToQueue(nowTime, event)
						}
					}
				}
				delete(eventQueues, nowTime)
			}
		}
	}
}

//无限次数执行
func DoTimer(delay int64, callback func()) uint64 {
	return Do(delay, 0, callback)
}

//延时处理
func SetTimeOut(delay int64, callback func()) uint64 {
	return Do(delay, 1, callback)
}

//移除一个定时器
func Remove(timerID uint64) {
	if _, exists := events[timerID]; exists {
		delete(events, timerID)
	}
}

//时间间隔，执行次数，回调函数
func Do(delay int64, repeatCount int, callback func()) uint64 {
	//最小单位1秒
	if delay < 1 {
		callback()
		return 0
	}

	if repeatCount < 0 {
		repeatCount = 0
	}

	event := &timerEvent{
		ID:          atomic.AddUint64(&maxSessionId, 1),
		Func:        callback,
		Delay:       delay,
		RepeatCount: repeatCount,
	}
	events[event.ID] = event

	addToQueue(time.Now().Unix(), event)

	return event.ID
}

//添加到执行队列
func addToQueue(nowTime int64, event *timerEvent) {
	dotime := nowTime + event.Delay
	events, exists := eventQueues[dotime]
	if !exists {
		events = list.New()
	}
	events.PushBack(event.ID)
	eventQueues[dotime] = events
}
