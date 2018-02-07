BINARY = pingback
SOURCE := *.go
BUILD_DIR=$(shell pwd)/bin

all: clean linux darwin windows

linux: 
	GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/pingback-linux-386 .
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/pingback-linux-amd64 .
	GOOS=linux GOARCH=arm GOARM=6 go build -o ${BUILD_DIR}/pingback-linux-arm .

windows:
	GOOS=windows GOARCH=386 go build -o ${BUILD_DIR}/pingback-winows-386.exe .
	GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/pingback-windows-amd64.exe .

darwin:
	GOOS=darwin GOARCH=386 go build -o ${BUILD_DIR}/pingback-darwin-386 .
	GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/pingback-darwin-amd64 .

clean:
	rm -rf bin/*

container: linux
	docker build -t rogierlommers/pingback-server .
	docker push rogierlommers/pingback-server:latest
