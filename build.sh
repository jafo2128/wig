#!/bin/bash

protoc --go_out=. *.proto
go build
cp -f msgs.proto www\
