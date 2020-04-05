SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
GO111MODULE=auto
go build -o ./selecttest  ./go-epoll/com/zzz/main/selecttest.go
pause
