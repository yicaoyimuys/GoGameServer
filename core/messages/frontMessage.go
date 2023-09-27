package messages

import (
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/sessions"
)

func FontReceive(session *sessions.FrontSession, msgBody []byte) {
	WARN("FontReceive自己实现一下吧，实现分发逻辑")
}
