package messages

import (
	. "GoGameServer/core/libs"
	"GoGameServer/core/libs/sessions"
)

func FontReceive(session *sessions.FrontSession, msgBody []byte) {
	WARN("FontReceive自己实现一下吧，实现分发逻辑")
}
