#!/bin/sh

node src/connectorServer/tools/proto/ClearMsg.js -n connectorServer

node src/connectorServer/tools/proto/CreateMessage.js -n connectorServer -p gameProto
node src/connectorServer/tools/proto/CreateMessage.js -n connectorServer -p systemProto