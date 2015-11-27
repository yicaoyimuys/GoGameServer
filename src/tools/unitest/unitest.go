package unitest

import (
	"testing"
	. "tools"
)

func Pass(t *testing.T, condition bool) bool {
	if condition {
		return true
	}
	ERR("[Unitest NOT PSSS]", "")
	t.FailNow()
	return false
}

func NotError(t *testing.T, err error) bool {
	if err == nil {
		return true
	}
	ERR("[Unitest ERROR]", err.Error())
	t.FailNow()
	return false
}
