#!/bin/sh

protoc --go_out=plugins=grpc:. src/core/libs/grpc/ipc/*.proto