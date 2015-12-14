#!/bin/sh 
function func(){
	killall -9 $1
	
	killall -0 $1
	while [ $? -ne 1 ]; do
		sleep 1
		killall -0 $1
	done
}

if [ $# -eq 0 ]
	then
		func GateServer
		sleep 1
		func LoginServer
		sleep 1
		func GameServer
		sleep 1
		func WorldServer
		sleep 1
		func DBServer
		sleep 1
		func LogServer
		sleep 1
	else
		func $1
fi