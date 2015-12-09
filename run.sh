#!/bin/sh 

sh stop.sh

$GOGAMESERVER_PATH/bin/LogServer &
sleep 1
$GOGAMESERVER_PATH/bin/DBServer &
sleep 1
$GOGAMESERVER_PATH/bin/GateServer &
sleep 1
$GOGAMESERVER_PATH/bin/LoginServer &
sleep 1
$GOGAMESERVER_PATH/bin/WorldServer &
sleep 1
$GOGAMESERVER_PATH/bin/GameServer -s=1 &
sleep 1
$GOGAMESERVER_PATH/bin/GameServer -s=2 &
sleep 1
$GOGAMESERVER_PATH/bin/GameServer -s=3
sleep 1
