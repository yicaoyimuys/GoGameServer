#! /bin/bash

sh stop_production.sh

pm2 start ./release/connectorServer -f --name bg_go_connector1 -- -e wxGame -g 1 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector2 -- -e wxGame -g 2 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector3 -- -e wxGame -g 3 &
sleep 1
pm2 start ./release/connectorServer -f --name bg_go_connector4 -- -e wxGame -g 4
