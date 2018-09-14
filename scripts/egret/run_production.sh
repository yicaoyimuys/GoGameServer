#! /bin/bash

sh stop_production.sh

pm2 start ./release/connectorServer -f --name bg_go_connector1 -- -e egret -g 1 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector2 -- -e egret -g 2 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector3 -- -e egret -g 3 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector4 -- -e egret -g 4
