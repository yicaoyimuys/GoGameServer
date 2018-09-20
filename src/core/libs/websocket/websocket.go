package websocket

import (
	. "core/libs"
	"core/libs/guid"
	"core/libs/stack"
	"core/sessions"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SessionMsgHandle func(session *sessions.FrontSession, msgBody []byte)
type SessionCloseHandle func(session *sessions.FrontSession)

type Server struct {
	port string
	guid *guid.Guid

	useSSL bool
	tslCrt string
	tslKey string

	sessionMsgHandle   SessionMsgHandle
	sessionCloseHandle SessionCloseHandle
}

func NewServer(port string, serviceId int) *Server {
	server := &Server{
		port:   port,
		guid:   guid.NewGuid(uint16(serviceId)),
		useSSL: false,
	}
	return server
}

func (this *Server) SetTLS(tslCrt string, tslKey string) {
	this.useSSL = true
	this.tslCrt = tslCrt
	this.tslKey = tslKey
}

func (this *Server) SetSessionMsgHandle(handle SessionMsgHandle) {
	this.sessionMsgHandle = handle
}

func (this *Server) SetSessionCloseHandle(handle SessionCloseHandle) {
	this.sessionCloseHandle = handle
}

func (this *Server) Start() {
	INFO("front start webSocket...", this.port)

	go func() {
		http.HandleFunc("/", this.wsHandler)
		var err error
		if this.useSSL {
			err = http.ListenAndServeTLS("0.0.0.0:"+this.port, this.tslCrt, this.tslKey, nil)
		} else {
			err = http.ListenAndServe("0.0.0.0:"+this.port, nil)
		}
		CheckError(err)
	}()
}

func (this *Server) StartPing() {
	overTime := 15
	sessions.FrontSessionOpenPing(int64(overTime))
	INFO("Session超时时间设置", overTime)
}

func (this *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ERR("wsHandler: ", err)
		return
	}

	//捕获异常
	defer stack.PrintPanicStackError()

	//Session创建
	sessionId := this.guid.NewID()
	sessionCodec := sessions.NewFrontCodec(conn)
	session := sessions.NewFontSession(sessionId, sessionCodec)
	this.addFontSession(session)
}

func (this *Server) addFontSession(session *sessions.FrontSession) {
	sessions.AddFrontSession(session)
	if this.sessionMsgHandle != nil {
		session.SetMsgHandle(this.sessionMsgHandle)
	}
	session.AddCloseCallback(nil, "webSocket.FrontSessionOffline", func() {
		if this.sessionCloseHandle != nil {
			this.sessionCloseHandle(session)
		}
		//DEBUG("session count: ", sessions.FrontSessionLen())
	})
	//DEBUG("session count: ", sessions.FrontSessionLen())

	defer session.Close()
	for {
		msg, err := session.Receive()
		if err != nil || msg == nil {
			break
		}
	}
}
