package random

import (
	"math/rand"
	"time"
)

func RandomInt31n(n int32) int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int31n(n)
}
