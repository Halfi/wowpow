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
	@go test ./...  -coverprofile=coverage.txt -covermode=atomic

.PHONY: coverage
coverage: test
	@go tool cover -func=coverage.txt

.PHONY: start
start:
	docker-compose -f deployments/docker-compose.yaml up
