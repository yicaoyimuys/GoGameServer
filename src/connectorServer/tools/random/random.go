package random

import (
	"math"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt31n(n int32) int32 {
	return r.Int31n(n)
}

// 随机int [0,n)
func RandIntn(n int) int {
	return r.Intn(n)
}

// 随机int [min,max)
func RandIntRange(min int, max int) int {
	if min >= max {
		return max
	}
	return r.Intn(max-min) + min
}

// 随机float64 [0.0,1.0)
func RandFloat64() float64 {
	return r.Float64()
}

// 随机数组中一个元素
func RandArray(arr []interface{}) interface{} {
	var index = math.Floor(r.Float64() * float64(len(arr)))
	return arr[int(index)]
}
