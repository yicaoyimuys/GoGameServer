#! /bin/bash

sh stop_production.sh

pm2 start ./release/connectorServer -f --name bg_go_connector1 -- -e wanba -g 1 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector2 -- -e wanba -g 2 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector3 -- -e wanba -g 3 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector4 -- -e wanba -g 4
