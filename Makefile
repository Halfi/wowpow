.PHONY: install
install:
	@go mod tidy
	@go install \
	  github.com/golang/mock/mockgen \
	  google.golang.org/protobuf/cmd/protoc-gen-go

.PHONY: gen
gen: install
	@go generate ./...

.PHONY: genProto
genProto: gen
	@cd api/proto && buf generate

.PHONY: test
test: gen
	@go test $$(go list ./... | grep -v /mock) -coverprofile=coverage.txt -covermode=atomic -timeout 60s

.PHONY: coverage
coverage: test
	@go tool cover -func=coverage.txt

.PHONY: lint
lint:
	@golangci-lint --timeout 10m0s run

.PHONY: start
start:
	docker-compose -f deployments/docker-compose.yaml up
