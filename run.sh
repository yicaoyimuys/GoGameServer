#!/bin/sh 

sh stop.sh

$GOGAMESERVER_PATH/bin/GameServer -s=1 &
$GOGAMESERVER_PATH/bin/GameServer -s=2 &
$GOGAMESERVER_PATH/bin/GameServer -s=3 &
sleep 1
$GOGAMESERVER_PATH/bin/GateServer
sleep 1
