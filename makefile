.PHONY: prepare #build test

$(eval $(service):;@:)
check:
	@[ "${service}" ] || ( echo "\x1b[31;1mERROR: 'service' is not set\x1b[0m"; exit 1 )
	@if [ ! -d "cmd/$(service)" ]; then  echo "\x1b[31;1mERROR: service '$(service)' undefined\x1b[0m"; exit 1; fi

build: check
	@go build -o build/$(service) cmd/$(service)/*.go

proto: check
	@if [ ! -d "sdk/$(service)" ]; then echo "creating new proto files..." &&  mkdir sdk/$(service) && mkdir sdk/$(service)/proto; fi
	$(foreach proto_file, $(shell find internal/services/$(service)/api/proto -name '*.proto'),\
	protoc --proto_path=internal/services/$(service)/api/proto \
		-I./sdk/proto \
		--go_out=sdk/$(service)/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=sdk/$(service)/proto \
		--go-grpc_opt=paths=source_relative $(proto_file) \
		--grpc-gateway_out sdk/$(service)/proto \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--openapiv2_out ./cmd/docs/swagger \
                --openapiv2_opt logtostderr=true \
                --openapiv2_opt use_go_templates=true )


