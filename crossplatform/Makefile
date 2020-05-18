all: lnx32 lnx64 win32 win64 mac64
.PHONY: lnx32 lnx64 win32 win64 mac64

lnx32:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o crossplatform/x86/yocoin cmd/yocoin/*.go
	cd crossplatform && tar -czf yocoin-linux32.tar.gz x86/yocoin

lnx64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o crossplatform/x64/yocoin cmd/yocoin/*.go
	cd crossplatform && tar -czf yocoin-linux64.tar.gz x64/yocoin


win32:
	CGO_ENABLED=1 CC=i686-w64-mingw32-gcc-win32 GOOS=windows GOARCH=386 go build -o crossplatform/win32/yocoin.exe cmd/yocoin/*.go
	cd crossplatform && zip -r yocoin-win32.zip  win32/yocoin.exe

win64:
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o crossplatform/win64/yocoin.exe cmd/yocoin/*.go
	cd crossplatform && zip -r yocoin-win64.zip  win64/yocoin.exe

mac64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o crossplatform/mac64/yocoin cmd/yocoin/*.go
	cd crossplatform && tar -czf yocoin-mac64.tar.gz mac64/
