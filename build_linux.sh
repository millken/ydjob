export GOPATH=$(cd "$(dirname "$0")"; pwd)
go build  -ldflags "-s"