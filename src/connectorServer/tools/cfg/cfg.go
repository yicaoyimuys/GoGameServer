package cfg

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ROOT string
)

func init() {
	initRootPath()
}

func initRootPath() {
	curFilename := os.Args[0]
	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}

	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}

	ROOT = filepath.Dir(filepath.Dir(binaryPath))
}
