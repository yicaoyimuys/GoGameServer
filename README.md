GoGameServer
===============

使用go搭建的一个游戏服务器项目  
1：consul作为服务注册/发现中心，便于服务的动态扩容  
2：vendor管理第三方库  
3：与客户端的通信协议使用protobuf  
4：服务之间支持基于TCP的ipc，rpc通信  
5：缓存使用redis  
6：数据库支持mysql、mongodb


结构图
===============


![image](GoGameServer.png)



Makefile脚本说明
===============

	make: 代码编译
	make clean: 清理编译文件
	make fmt: 代码格式化
	make proto: 生成protobuf的go文件


启动说明
===============

	执行make进行代码编译
	启动consul
	启动mysql(config/local/mysql.json)
	启动redis(config/local/redis.json)
	执行sh run.sh启动服务器
    执行sh test.sh启动测试服务器