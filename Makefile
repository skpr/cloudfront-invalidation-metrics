# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run revive to lint the code
lint:
	revive -config revive.toml -set_exit_status ./...

test: fmt vet
	go test ./...

