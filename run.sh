#!/bin/sh 

sh stop.sh

./bin/connectorServer -e local -s 1 &
sleep 1
./bin/connectorServer -e local -s 2
