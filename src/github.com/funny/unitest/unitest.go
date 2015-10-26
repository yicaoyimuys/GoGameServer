package unitest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var (
	passRegexp, _     = regexp.Compile(`^\s*unitest\.Pass\s*\(\s*[^,]+\s*,\s*(.+)\s*\)\s*$`)
	notErrorRegexp, _ = regexp.Compile(`^\s*unitest\.NotError\s*\(\s*[^,]+\s*,\s*(.+)\s*\)\s*$`)
)

func Pass(t *testing.T, condition bool) bool {
	if condition {
		return true
	}
	log("[NOT PSSS]", passRegexp, "")
	t.FailNow()
	return false
}

func NotError(t *testing.T, err error) bool {
	if err == nil {
		return true
	}
	log("[ERROR]", notErrorRegexp, err.Error())
	t.FailNow()
	return false
}

func log(title string, regex *regexp.Regexp, val string) {
	if _, file, line, ok := runtime.Caller(2); ok {
		if data, err := ioutil.ReadFile(file); err == nil {
			// Truncate file name at last file name separator.
			if index := strings.LastIndex(file, "/"); index >= 0 {
				file = file[index+1:]
			} else if index = strings.LastIndex(file, "\\"); index >= 0 {
				file = file[index+1:]
			}
			lines := bytes.Split(data, []byte{'\n'})
			cond := regex.FindAllSubmatch(lines[line-1], 1)
			if len(cond) > 0 && len(cond[0]) > 1 {
				if val == "" {
					fmt.Fprintf(os.Stderr, "\t%s %s:%d: %s\n", title, file, line, cond[0][1])
				} else {
					fmt.Fprintf(os.Stderr, "\t%s %s:%d: %s: %s\n", title, file, line, cond[0][1], val)
				}
			}
		}
	}
}
