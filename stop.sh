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
		func LoginServer
		func GameServer
		func WorldServer
		func DBServer
		func LogServer
	else
		func $1
fi