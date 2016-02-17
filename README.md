## Wig (go-web-irc)

A very simple self hostable FCGI irc client. 

#### Setup 

```git clone https://github.com/v0l/wig.git```

```
cd wig
protoc --go_out=. .\proto\*.proto
go build
```