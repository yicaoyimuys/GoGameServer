#!/bin/sh

node src/tools/proto/ClearMsg.js -n connectorServer

node src/tools/proto/CreateMessage.js -n connectorServer -p gameProto
node src/tools/proto/CreateMessage.js -n connectorServer -p systemProto