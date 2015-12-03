package systemProto

const (
	//连接DB服务器
	ID_System_ConnectDBServerC2S = 10001
	ID_System_ConnectDBServerS2C = 10002

	//连接Transfer服务器
	ID_System_ConnectTransferServerC2S = 10003
	ID_System_ConnectTransferServerS2C = 10004

	//连接World服务器
	ID_System_ConnectWorldServerC2S = 10005
	ID_System_ConnectWorldServerS2C = 10006

	//客户端Session状态使用
	ID_System_ClientSessionOnlineC2S  = 10007
	ID_System_ClientSessionOfflineC2S = 10009

	//客户端登录状态使用
	ID_System_ClientLoginSuccessC2S = 10011
	ID_System_ClientLoginSuccessS2C = 10012
)
