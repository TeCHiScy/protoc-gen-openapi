
all:
	go generate ./...
	go get ./...
	go install ./...

test:
	# since some tests call separately-built binaries, clear the cache to ensure all get run
	go clean -testcache
	go test ./... -v