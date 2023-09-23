

go env -w GOARCH=386
go env -w CGO_ENABLED=1
go build -buildmode=c-shared -o main.dll
go env -u GOARCH
go env -u CGO_ENABLED

del /q          ..\main.dll
move main.dll ..\


