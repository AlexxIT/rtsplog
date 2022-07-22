@ECHO OFF

@SET GOOS=windows
@SET GOARCH=amd64
@SET FILENAME=rtsplog_amd64.exe
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=amd64
@SET FILENAME=rtsplog_amd64
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=386
@SET FILENAME=rtsplog_i386
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=arm
@SET GOARM=7
@SET FILENAME=rtsplog_armv7
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=arm64
@SET FILENAME=rtsplog_aarch64
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=darwin
@SET GOARCH=amd64
@SET FILENAME=rtsplog_darwin
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%
