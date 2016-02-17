@echo off

protoc --go_out=. *.proto
go build
copy /Y .\msgs.proto www\msgs.proto
echo Done.