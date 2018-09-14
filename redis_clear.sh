#!/bin/sh 
redis-cli -p 20036 -a yangsong flushdb &
sleep 1
redis-cli -p 20036 -a yangsong flushall &
sleep 1