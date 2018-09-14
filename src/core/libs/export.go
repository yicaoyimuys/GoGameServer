package libs

import (
	"fmt"
	"strings"
)

import (
	"core/libs/logger"
	"github.com/astaxie/beego/logs"
	"net"
	"reflect"
	"strconv"
)

var (
	loggerPrefix string
	log          = logs.GetBeeLogger()
)

func init() {
}

//------------------------------------------------ 严重程度由高到低
func ERR(v ...interface{}) {
	log.Error(getLogMsg(v))
}

func WARN(v ...interface{}) {
	log.Warn(getLogMsg(v))
}

func INFO(v ...interface{}) {
	log.Info(getLogMsg(v))
}

func NOTICE(v ...interface{}) {
	log.Notice(getLogMsg(v))
}

func DEBUG(v ...interface{}) {
	log.Debug(getLogMsg(v))
}

func getLogMsg(v ...interface{}) string {
	return "[" + loggerPrefix + "] " + strings.TrimRight(fmt.Sprintln(v...), "\n")
}

func SetLogDebug(debug bool) {
	if debug {
		log.SetLevel(logs.LevelDebug)
	} else {
		log.SetLevel(logs.LevelInfo)
	}
}

func SetLogFile(fileName string, log_output string) {
	loggerPrefix = fileName
	logger.StartLogger(loggerPrefix+".log", log_output)
}

// 保持进程
func Run() {
	temp := make(chan int32, 10)
	for {
		select {
		case <-temp:
		}
	}
}

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

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IndexOf(array interface{}, value interface{}) int {
	arrType := reflect.TypeOf(array).Kind()
	if arrType != reflect.Slice && arrType != reflect.Array {
		return -1
	}

	arr := reflect.ValueOf(array)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == value {
			return i
		}
	}
	return -1
}

func InArray(array interface{}, value interface{}) bool {
	return IndexOf(array, value) != -1
}
