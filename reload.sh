#!/bin/bash
function func(){
	a=$(ps aux | grep $1 | awk '{print $2}')

	for i in $a
	do	
		#echo $i
		kill -1 $i
	done
}

func GateServer
func GameServer