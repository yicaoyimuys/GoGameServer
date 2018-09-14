#! /bin/bash

serverHostName=`hostname`
sidArr=($(tr "-" " " <<< $serverHostName))


start=${sidArr[0]}
sId=0
if [[ $start == "xm" ]]; then
	sId=${sidArr[3]}	
fi

sh stop_production.sh

pm2 start ./release/connectorServer -f --name bg_go_connector1 -- -e baiwan -g 1 -s $sId &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector2 -- -e baiwan -g 2 -s $sId &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector3 -- -e baiwan -g 3 -s $sId &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector4 -- -e baiwan -g 4 -s $sId
