package file

import (
	"os"
	. "tools"
)

//创建循环目录，如果存在则不执行任何操作
func CreateDir(dirPath string) error {
	dirPath = os.Getenv("GOGAMESERVER_PATH") + dirPath
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil{
		ERR("CreateDir: ", err)
	}
	return err
}

//打开文件，如果不存在则创建
func OpenFile(filePath string) *os.File {
	filePath = os.Getenv("GOGAMESERVER_PATH") + filePath
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil{
		ERR("OpenFile: ", err)
		return nil
	}
	return file
}