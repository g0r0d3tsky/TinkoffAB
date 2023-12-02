PHONY: generate
generate:
	protoc --go_out=./ --go_opt=paths=source_relative \
			--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
			# TODO: --validate_out="lang=go,paths=source_relative:." validation_service.service \
			service/service.service
