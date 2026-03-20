PROTO_DIR=proto/bookstore

.PHONY: proto run-gateway run-user run-order run-catalog run-payment

proto:
	@echo "Generating protobuf stubs..."
	protoc -I proto \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  proto/bookstore/common.proto \
	  proto/bookstore/user.proto \
	  proto/bookstore/catalog.proto \
	  proto/bookstore/order.proto \
	  proto/bookstore/payment.proto
	python -m grpc_tools.protoc -I proto \
	  --python_out=. --grpc_python_out=. \
	  proto/bookstore/common.proto \
	  proto/bookstore/user.proto \
	  proto/bookstore/catalog.proto \
	  proto/bookstore/order.proto \
	  proto/bookstore/payment.proto

run-gateway:
	cd gateway-go && go run ./cmd/gateway

run-user:
	cd user-service-go && go run ./cmd/user

run-order:
	cd order-service-go && go run ./cmd/order

run-catalog:
	cd catalog-service-py && python -m app.main

run-payment:
	cd payment-service-py && python -m app.main
