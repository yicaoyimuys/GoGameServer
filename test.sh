#!/bin/sh 

sh stop.sh servives/test/main
sh stop.sh servives/test/main_ws
go run servives/test/main.go -e local
# go run servives/test/main_ws.go -e local