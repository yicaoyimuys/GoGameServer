package common

import (
	"net"
	"reflect"
	"strconv"
	"time"
)

//获取当前毫秒时间戳
func UnixMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

//数字转成字符串
func NumToString(num interface{}) string {
	str := ""
	numType := reflect.TypeOf(num).Kind()
	switch numType {
	case reflect.Int8:
		str = strconv.FormatInt(int64(num.(int8)), 10)
	case reflect.Int16:
		str = strconv.FormatInt(int64(num.(int16)), 10)
	case reflect.Int32:
		str = strconv.FormatInt(int64(num.(int32)), 10)
	case reflect.Int64:
		str = strconv.FormatInt(int64(num.(int64)), 10)
	case reflect.Int:
		str = strconv.FormatInt(int64(num.(int)), 10)
	case reflect.Uint8:
		str = strconv.FormatUint(uint64(num.(uint8)), 10)
	case reflect.Uint16:
		str = strconv.FormatUint(uint64(num.(uint16)), 10)
	case reflect.Uint32:
		str = strconv.FormatUint(uint64(num.(uint32)), 10)
	case reflect.Uint64:
		str = strconv.FormatUint(uint64(num.(uint64)), 10)
	case reflect.Uint:
		str = strconv.FormatUint(uint64(num.(uint)), 10)
	case reflect.Float64:
		str = strconv.FormatFloat(num.(float64), 'f', -1, 64)
	}
	return str
}

//Float转成字符串
func FloatToString(num interface{}, fixedLen int) string {
	str := ""
	numType := reflect.TypeOf(num).Kind()
	switch numType {
	case reflect.Float32:
		str = strconv.FormatFloat(float64(num.(float32)), 'f', fixedLen, 64)
	case reflect.Float64:
		str = strconv.FormatFloat(float64(num.(float64)), 'f', fixedLen, 64)
	default:
		str = NumToString(num)
	}
	return str
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
