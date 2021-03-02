package guid

import (
	"GoGameServer/core/libs/logger"
	"sync"
	"time"
)

type Guid struct {
	serverId      uint16
	sequence      int32
	mx            sync.RWMutex
	lastTimestamp int64
}

func (this *Guid) NewID() uint64 {
	this.mx.Lock()
	defer this.mx.Unlock()

	if this.serverId > 4095 {
		logger.Error("server_id超出最大值")
		return 0
	}

	timestamp := time.Now().Unix()
	if timestamp < this.lastTimestamp {
		logger.Error("请调整服务器时间!")
		return 0
	}

	if timestamp == this.lastTimestamp {
		// 当前毫秒内，则+1
		this.sequence += 1
		if this.sequence > 4095 {
			// 当前毫秒内计数满了，则等待下一秒
			this.sequence = 0
			for {
				timestamp = time.Now().Unix()
				if timestamp > this.lastTimestamp {
					break
				}
			}
		}
	} else {
		this.sequence = 0
	}
	this.lastTimestamp = timestamp

	//40(毫秒) + 12(serverID) + 12(重复累加)
	return uint64(timestamp<<40) | (uint64(this.serverId) << 12) | uint64(this.sequence)
}

func NewGuid(server_id uint16) *Guid {
	return &Guid{
		serverId:      server_id,
		sequence:      0,
		lastTimestamp: -1,
	}
}
