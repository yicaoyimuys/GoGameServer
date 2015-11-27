package cfg

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	_map  map[string]string
	_lock sync.Mutex
)

func init() {
	Reload()
}

func Get() map[string]string {
	_lock.Lock()
	defer _lock.Unlock()
	return _map
}

func Reload() {
	var baseConfigPath string = os.Getenv("GOGAMESERVER_PATH") + "/config_base.ini"
	var serverConfigPath string = os.Getenv("GOGAMESERVER_PATH") + "/config_server.ini"

	_lock.Lock()
	_map = make(map[string]string)
	_load_config(baseConfigPath)
	_load_config(serverConfigPath)
	_lock.Unlock()
}

func _load_config(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return
	}

	re := regexp.MustCompile(`[\t ]*([0-9A-Za-z_]+)[\t ]*=[\t ]*([^\t\n\f\r# ]+)[\t #]*`)

	// using scanner to read config file
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// expression match
		slice := re.FindStringSubmatch(line)

		if slice != nil {
			_map[slice[1]] = slice[2]
			//			log.Println(slice[1], "=", slice[2])
		}
	}

	return
}

func GetValue(key string) string {
	config := Get()
	return config[key]
}

func GetUint16(key string) uint16 {
	result, _ := strconv.Atoi(GetValue(key))
	return uint16(result)
}

func GetInt(key string) int {
	result, _ := strconv.Atoi(GetValue(key))
	return result
}
