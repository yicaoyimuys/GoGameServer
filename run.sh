#!/bin/sh 

sh stop.sh

./bin/connectorServer -e development -g 1 &
sleep 1
./bin/connectorServer -e development -g 2
