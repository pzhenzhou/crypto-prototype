BINARY_NAME=crypto-http

go=GO111MODULE=on go
flags="-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse --short HEAD`"

build-linux64:
	GOARCH=amd64 GOOS=linux ${go} build -ldflags ${flags} -x -o bin/${BINARY_NAME}-linux main.go

build-darwin:
	GOARCH=amd64 GOOS=darwin ${go} build -ldflags ${flags} -x -o bin/${BINARY_NAME}-darwin main.go

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 ${go} build -ldflags ${flags} -x -o bin/${BINARY_NAME}-arm64 ./cmd/main.go

clean:
	go clean
	rm bin/${BINARY_NAME}-*

