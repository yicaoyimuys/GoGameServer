#!/bin/sh 

DIR=$(pwd)
cd $DIR/core/protos/gameProto && protoc --go_out=. gameProto.proto
cd $DIR/core/libs/grpc/ipc && protoc --go_out=plugins=grpc:. *.proto