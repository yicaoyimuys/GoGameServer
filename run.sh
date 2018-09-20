#!/bin/sh 

sh stop.sh

./bin/connector -e local -s 1 &
sleep 1
./bin/connector -e local -s 2
