#!/bin/sh 
function func(){
	ps -ef | grep $1 | grep -v grep | awk '{print $2}' | xargs kill -9
}

if [ $# -eq 0 ]
	then
		func servives/connector/main
		func servives/game/main
		func servives/login/main
		func servives/chat/main
		func servives/log/main
		func servives/api/main
		func exe/main
	else
		func $1
fi