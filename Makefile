kustomtree:
	go build -o $@ .

.PHONY: test
test:
	golangci-lint run
	$(MAKE) -C integration_test
