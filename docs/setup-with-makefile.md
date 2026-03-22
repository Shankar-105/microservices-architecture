# Setup With Makefile

This guide shows how to run the project locally using Makefile targets.

## Prerequisites
1. Go (1.22+ recommended)
2. Python (3.12+ recommended)
3. `protoc` available in PATH
4. `make`

## What Makefile Handles For You
When you use the provided Make targets, the setup steps are already automated:
1. Creates a local Python virtual environment in `.venv`.
2. Installs Python requirements for catalog service.
3. Installs/uses Go proto plugins (`protoc-gen-go`, `protoc-gen-go-grpc`).
4. Generates Go and Python stubs from proto contracts.
5. Starts all services with the right startup commands.

## 1) Generate gRPC/Proto Stubs
From repository root:

```bash
make proto
```

Generated output:
1. Go stubs in `generated-go/`
2. Python stubs in `generated-py/`

Tip:
1. If this is your first run on a machine, run `make install-go-plugins` once before `make proto`.

## 2) (Optional but Recommended) Run Setup Targets Explicitly
If you are a beginner and want to see setup happen step by step, run these in order:

```bash
make init-py
make install-py-deps
make install-go-plugins
make proto
```

What these do:
1. `make init-py`: creates `.venv` and upgrades `pip/setuptools/wheel`.
2. `make install-py-deps`: installs catalog service Python packages from `requirements.txt`.
3. `make install-go-plugins`: installs Go plugins used by `protoc`.
4. `make proto`: generates all stubs.

## 3) Start Services

```bash
make run-all
```

Services started:
1. Gateway HTTP on `:8080`
2. User service gRPC on `:50051`
3. Order service gRPC on `:50052`
4. Catalog service gRPC on `:50053` and HTTP on `:8001`

## 4) Verify Ports

```bash
make status-ports
```

## 5) Test Routes
1. `GET http://localhost:8080/healthz`
2. `GET http://localhost:8080/getallbooks`
3. `POST http://localhost:8080/createorder`
	- JSON body: `{"user_id":"u-100","book_id":"b-100","quantity":1}`
4. `GET http://localhost:8080/getorders?user_id=u-100`

## 6) Stop Services

```bash
make stop-all
```

## Useful Individual Targets
1. `make run-gateway`
2. `make run-user`
3. `make run-order`
4. `make run-catalog`

## Common Beginner Issues
1. `protoc not found`:
	- Install Protocol Buffers compiler and ensure `protoc` is in PATH.
2. Python venv issues:
	- Delete `.venv` and rerun `make init-py`.
3. Go plugin missing errors:
	- Run `make install-go-plugins`.
