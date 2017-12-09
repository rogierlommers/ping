source := *.go

release:
	# create directories
	mkdir -p ./bin/windows64
	mkdir -p ./bin/macos64
	mkdir -p ./bin/linux64
	mkdir -p ./bin/linux32

	# build binaries
	#GOOS=windows GOARCH=amd64 go build -o ./bin/windows64/pingback ${source}
	#GOOS=darwin GOARCH=amd64 go build -o ./bin/macos64/pingback ${source}
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux64/pingback ${source}
	GOOS=linux GOARCH=386 go build -o ./bin/linux32/pingback ${source}
