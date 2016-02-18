#!/bin/bash

pwd
mkdir build

cd build

wget https://github.com/google/protobuf/releases/download/v2.6.1/protobuf-2.6.1.tar.gz
tar xfv protobuf-2.6.1.tar.gz
cd protobuf-2.6.1
./autogen.sh
./configure
make -j4
cd ..

pwd
./protobuf-2.6.1/src/protoc --go_out=.. --proto_path=.. ../*.proto
go build ../
cp -f ../msgs.proto ../www/

rm -rf protobuf-2.6.1*
