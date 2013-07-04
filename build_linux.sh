#!/bin/sh
#go get -v
export GOPATH=$(cd "$(dirname "$0")"; pwd)
echo $GOPATH
REVISION=`git rev-parse --short=5 HEAD`
echo $REVISION > REVISION
cd src 
go build -ldflags "-s -X main.gitVersion $REVISION" -v #会在当前目录
go install ydjob #这样编译的会在bin下
cd -
