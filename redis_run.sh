#!/bin/sh 

#需要配置Redis环境变量
#export PATH=$PATH:/Users/egret/Documents/redis-3.0.5/src
redis-server ./config/redis.conf &
sleep 1
ps -ef | grep redis