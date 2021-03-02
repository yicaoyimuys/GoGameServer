#!/bin/sh 

sh stop.sh

sleep 1
go run servives/api/main.go -e local -s 1 &

sleep 1
go run servives/log/main.go -e local -s 1 &

sleep 1
go run servives/game/main.go -e local -s 1 &
sleep 1
go run servives/game/main.go -e local -s 2 &

sleep 1
go run servives/login/main.go -e local -s 1 &

sleep 1
go run servives/chat/main.go -e local -s 1 &

sleep 1
go run servives/connector/main.go -e local -s 1 &
sleep 1
go run servives/connector/main.go -e local -s 2