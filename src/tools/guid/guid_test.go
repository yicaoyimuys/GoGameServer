package guid

import "code.google.com/p/goprotobuf/proto"
import (
	//	"bytes"
	//	"encoding/binary"
	"fmt"
	//	"strconv"
	"sync"
	"testing"
	//	. "tools"
)

var mx sync.Mutex
var test map[uint64]int

func TestNewID(t *testing.T) {
	guid := NewGuid()

	test = make(map[uint64]int)
	for i := 0; i < 1; i++ {
		go func() {
			//			id := guid.NewID(uint16(RandomInt31n(4095)))
			id := guid.NewID(1)
			//			fmt.Println(id)
			addID(id)
		}()
	}

	//	fmt.Println(int8(byte(128)))
	//	fmt.Println(uint8(byte(128)))

	//	var a []byte = make([]byte, 8)
	//	a[0] = 128
	//	a[1] = 128
	//	a[2] = 128
	//	a[3] = 128
	//	a[4] = 128
	//	a[5] = 128
	//	a[6] = 128
	//	a[7] = 128

	//	fmt.Println(binary.BigEndian.Uint64(a))

	//	b_buf := bytes.NewBuffer(a)
	//	var x int64
	//	binary.Read(b_buf, binary.BigEndian, &x)
	//	fmt.Println(x)

	fmt.Println(proto.EncodeVarint(322590098456576))
}

func addID(id uint64) {
	mx.Lock()
	defer mx.Unlock()

	if _, exists := test[id]; exists {
		fmt.Println("what?", id)
		return
	}
	test[id] = 0
}
