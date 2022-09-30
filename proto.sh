#!/bin/sh 

# go get -u google.golang.org/protobuf/cmd/protoc-gen-go
# go install google.golang.org/protobuf/cmd/protoc-gen-go
# go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

DIR=$(pwd)
cd $DIR/servives/public/gameProto && protoc --go_out=. gameProto.proto
cd $DIR/core/libs/grpc/ipc && protoc --go_out=. --go-grpc_out=. ipc.proto