package common

import (
	"net"
	"time"
)

//获取当前毫秒时间戳
func UnixMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

//获取本机Ip
func GetLocalIp() string {
	ipAddr := "localhost"
	addrSlice, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrSlice {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if nil != ipnet.IP.To4() {
					ipAddr = ipnet.IP.String()
					break
				}
			}
		}
	}
	return ipAddr
}

//三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
