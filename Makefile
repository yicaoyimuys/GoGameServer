.PHONY: .FORCE
GO=go
DOT=dot
GOYACC=$(GO) tool yacc

SRC_DIR = ./src
PROTO_INSTALL_FILE_DIR = ./src/code.google.com/p/goprotobuf/

all:
	$(GO) install GateServer
	$(GO) install WorldServer
	$(GO) install GameServer
	$(GO) install DBServer
	$(GO) install LoginServer
	
local:
	$(GO) install LocalServer

clean:
	rm -rf bin pkg
	rm -rf dump.rdb
	rm -rf logs/db.log logs/game.log logs/gateway.log logs/local.log logs/login.log
 
fmt:
	$(GO) fmt $(SRC_DIR)/...


#安装proto
install_proto:
	make -C $(PROTO_INSTALL_FILE_DIR)
	
#需要先install_proto，然后将bin目录加入到环境变量，protoc才可使用
create_proto:
	cd $(SRC_DIR)/protos/systemProto && protoc --go_out=. systemProto.proto
	cd $(SRC_DIR)/protos/dbProto && protoc --go_out=. dbProto.proto
	cd $(SRC_DIR)/protos/gameProto && protoc --go_out=. gameProto.proto
	
#交叉编译：
#首先进入go源码目录
#cd /usr/local/go/src/
#执行sudo GOOS=linux GOARCH=amd64 ./make.bash  生成linux下的编译文件pkg
#以上步骤在go1.5中应该已不再需要，但还未测试
publish_linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o GateServer_linux GateServer
	GOOS=linux GOARCH=amd64 $(GO) build -o GameServer_linux GameServer
	
publish_windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o GateServer_windows.exe GateServer
	GOOS=windows GOARCH=amd64 $(GO) build -o GameServer_windows.exe GameServer
	
publish_mac:
	GOOS=darwin GOARCH=amd64 $(GO) build -o GateServer_darwin GateServer
	GOOS=darwin GOARCH=amd64 $(GO) build -o GameServer_darwin GameServer
	