SET CGO_ENABLED=0
SET GOARCH=amd64
SET GOOS=windows
go build -ldflags "-s -w"
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w"