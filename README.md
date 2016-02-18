## Wig (web-irc-go) 
![build](https://travis-ci.org/v0l/wig.svg?branch=master)

A very simple self hostable irc client built with websockets and ProtoBuf.

Most of the implementation is done in JS the go server only passes messages over and back and handles logins.

### Setup 

ProtoC: [download](https://developers.google.com/protocol-buffers/docs/downloads)

```
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/golang/protobuf/proto
go get github.com/gorilla/websocket
```

```
git clone https://github.com/v0l/wig.git
cd wig
./build
./wig
```

### Troubleshooting

 * ```--go_out: protoc-gen-go: The system cannot find the file specified.```
   * Make sure ```%GOPATH%\bin``` is added to your path (or copy ```%GOPATH%\bin\protoc-gen-go``` to your path)