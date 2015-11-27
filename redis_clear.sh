#!/bin/sh 
redis-cli flushdb &
sleep 1
redis-cli flushall &
sleep 1