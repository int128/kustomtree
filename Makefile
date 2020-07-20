kustomtree:
	go build -o $@ .

.PHONY: test
test:
	golangci-lint run
	go test -v ./...
	$(MAKE) -C integration_test
