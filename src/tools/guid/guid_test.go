package guid

import (
	"fmt"
	"testing"
)

func TestNewID(t *testing.T) {
	guid := NewGuid()

	for i := 0; i < 6002120; i++ {
		go func() {
			fmt.Println(guid.NewID(1, 1))
		}()
	}
}
