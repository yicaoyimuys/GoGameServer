#!/bin/sh 

sh stop.sh

./bin/connectorServer -e local -g 1 &
sleep 1
./bin/connectorServer -e local -g 2
