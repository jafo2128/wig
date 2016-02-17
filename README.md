## Wig (go-web-irc)

A very simple self hostable FCGI irc client. 

#### Setup 

Protoc: https://developers.google.com/protocol-buffers/docs/downloads

```
go get -u github.com/golang/protobuf/protoc-gen-go
git clone https://github.com/v0l/wig.git
cd wig
protoc --go_out=. .\proto\*.proto
go build
./wig
```