# Setup With Makefile

This guide shows how to run the project locally using Makefile targets.

## Prerequisites
1. Go (1.22+ recommended)
2. Python (3.12+ recommended)
3. `protoc` available in PATH
4. `make`

## 1) Generate gRPC/Proto Stubs
From repository root:

```bash
make proto
```

Generated output:
1. Go stubs in `generated-go/`
2. Python stubs in `generated-py/`

## 2) Start Services

```bash
make run-all
```

Services started:
1. Gateway HTTP on `:8080`
2. User service gRPC on `:50051`
3. Order service gRPC on `:50052`
4. Catalog service gRPC on `:50053` and HTTP on `:8001`

## 3) Verify Ports

```bash
make status-ports
```

## 4) Test Routes
1. `GET http://localhost:8080/healthz`
2. `GET http://localhost:8080/getallbooks`
3. `POST http://localhost:8080/createorder`
	- JSON body: `{"user_id":"u-100","book_id":"b-100","quantity":1}`
4. `GET http://localhost:8080/getorders?user_id=u-100`

## 5) Stop Services

```bash
make stop-all
```

## Useful Individual Targets
1. `make run-gateway`
2. `make run-user`
3. `make run-order`
4. `make run-catalog`
