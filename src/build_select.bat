SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
GO111MODULE=auto
go build -o ./select  ./go-epoll/com/zzz/main/select.go
pause
