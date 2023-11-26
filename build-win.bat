go mod tidy
go generate
go build -ldflags "-s -w"
