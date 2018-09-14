#!/bin/sh

protoc --go_out=plugins=grpc:. src/tools/grpc/ipc/*.proto