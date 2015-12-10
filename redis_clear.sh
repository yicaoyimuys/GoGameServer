#!/bin/sh 
redis-cli -p 6379 -a yangsong flushdb &
sleep 1
redis-cli -p 6379 -a yangsong flushall &
sleep 1