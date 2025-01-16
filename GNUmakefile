default: lint testacc generate

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m


# run golanci-lint
lint:
	golangci-lint run

# run go generate
generate:
	go generate ./...


.PHONY: testacc lint
