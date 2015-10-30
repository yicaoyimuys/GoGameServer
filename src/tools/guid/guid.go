package guid

import (
	"encoding/binary"
	"sync"
	"time"
)

type Guid struct {
	id uint16
	mx sync.RWMutex
}

func (this *Guid) NewID(platform_id uint16, server_id uint16) uint64 {
	this.mx.Lock()
	defer this.mx.Unlock()

	this.id += 1

	var b []byte = make([]byte, 8)
	binary.BigEndian.PutUint32(b[0:4], uint32(time.Now().Unix()))
	binary.BigEndian.PutUint16(b[4:6], uint16(platform_id*100+server_id))
	binary.BigEndian.PutUint16(b[6:8], this.id)

	if this.id == 65535 {
		this.id = 0
	}

	return binary.BigEndian.Uint64(b)
}

func NewGuid() *Guid {
	return &Guid{
		id: 0,
	}
}
