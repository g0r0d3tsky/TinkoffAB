PHONY: generate
generate:
	protoc --go_out=./ --go_opt=paths=source_relative \
			--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
			service/service.proto
			# TODO: --validate_out="lang=go,paths=source_relative:." validation_service.service \

build-docker-server:
	docker build -t trans-server -f Dockerfile.server .

build-docker-client:
	docker build -t trans-client -f Dockerfile.client .
