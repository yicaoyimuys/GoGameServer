package service

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/cast"
	. "github.com/yicaoyimuys/GoGameServer/core/libs"
	"github.com/yicaoyimuys/GoGameServer/core/libs/stack"
	"go.uber.org/zap"
)

func (this *Service) StartPProf(port int) {
	port = port + this.id
	go func() {
		defer stack.TryError()
		http.ListenAndServe(":"+cast.ToString(port), nil)
	}()
	INFO("PProf Start", zap.Int("Port", port))
}
