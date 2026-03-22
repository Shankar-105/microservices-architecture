# Setup With Docker

This guide shows how to run the project using Docker Compose.

## Prerequisites
1. Docker Desktop installed and running
2. Docker Compose v2 available (`docker compose`)

## 1) Build and Start Services
From repository root:

```bash
docker compose up --build
```

This starts:
1. `gateway-go`
2. `user-service-go`
3. `order-service-go`
4. `catalog-service-py`

## 2) Port Map
1. Gateway HTTP: `8080`
2. User gRPC: `50051`
3. Order gRPC: `50052`
4. Catalog gRPC: `50053`
5. Catalog HTTP: `8001`

## 3) Verify Health
1. Gateway health:
	- `GET http://localhost:8080/healthz`
2. Catalog health:
	- `GET http://localhost:8001/healthz`

## 4) Smoke Test API Routes
1. Get all books:
	- `GET http://localhost:8080/getallbooks`
2. Create order:
	- `POST http://localhost:8080/createorder`
	- Body:
	  - `{"user_id":"u-100","book_id":"b-100","quantity":1}`
3. Get orders:
	- `GET http://localhost:8080/getorders?user_id=u-100`

## 5) Stop Services

```bash
docker compose down
```

Optional cleanup of built images/volumes:

```bash
docker compose down --rmi local --volumes
```
