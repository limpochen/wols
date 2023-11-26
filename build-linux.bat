@echo off
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
echo on
go mod tidy
rem go generate
rem go build -ldflags "-s -w"
go build
@echo off
SET CGO_ENABLED=
SET GOOS=
SET GOARCH=
