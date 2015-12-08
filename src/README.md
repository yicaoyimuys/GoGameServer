<pre>
	GateServer -- 游戏网关服务器
	LoginServer -- 登录认证服务器
	GameServer -- 游戏服务器(处理快速逻辑)
	WorldServer -- 世界服务器(处理慢速逻辑)
	DBServer -- DB服务器
</pre>

<pre>
    dao -- DB访问操作模块
    global -- 全局使用的功能
	model -- 实体类
	module -- 各个模块
	protos -- 消息协议goprotobuf
	proxys -- 各个服务器间通信使用
	test -- 测试代码
	tools -- 工具类
</pre>

<pre>
	### 感谢:
	- github.com/funny/link -- 网络通信模块
	- github.com/go-sql-driver/mysql -- mysql-driver
	- github.com/hoisie/redis -- Redis访问操作模块
	- github.com/spaolacci/murmur3 -- MurmurHash算法
</pre>

<pre>
    #GO基本配置
    export GOROOT=/usr/local/go1.4.2
    export PATH=$PATH:$GOROOT/bin
    export GOPATH=$GOROOT

	#Redis配置
    export PATH=$PATH:/Users/egret/Documents/redis-3.0.5/src

    #protobuf配置
    安装protobuf

	GoGameServer使用配置
	export GOGAMESERVER_PATH=/Users/GoGameServer/
	export GOPATH=$GOPATH:$GOGAMESERVER_PATH
	export PATH=$PATH:$GOGAMESERVER_PATH/bin
</pre>
