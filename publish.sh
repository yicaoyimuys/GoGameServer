#! /bin/bash

make

make publish_linux

cp release/connectorServer ../__Publish/bgServerGo_Connector/release
cp config/log.json ../__Publish/bgServerGo_Connector/config
cp config/redis.json ../__Publish/bgServerGo_Connector/config
cp config/server.json ../__Publish/bgServerGo_Connector/config
cp -r scripts ../__Publish/bgServerGo_Connector