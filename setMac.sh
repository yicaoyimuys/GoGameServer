#!/bin/sh 

sudo sysctl -w kern.maxfiles=1048600
sudo sysctl -w kern.maxfilesperproc=1048576
sudo ulimit -n 65536
sudo sysctl -w kern.ipc.somaxconn=8192