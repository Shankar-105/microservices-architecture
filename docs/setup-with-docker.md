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

If images are already built and you only want to start containers:

```bash
docker compose up
```

If you want to run in detached mode (background):

```bash
docker compose up -d
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

## Essential Docker Commands You Should Know

These are the most useful docker commands you need to know while working with docker.

1. Build images only:

```bash
docker compose build
```

2. See running containers for this compose project:

```bash
docker compose ps
```

3. Watch logs of all services:

```bash
docker compose logs -f
```

4. Watch logs for one service:

```bash
docker compose logs -f gateway-go
```

5. Restart a single service:

```bash
docker compose restart gateway-go
```

6. Rebuild one service after code changes:

```bash
docker compose build gateway-go
docker compose up -d gateway-go
```

7. Open shell or bash in a running container:

```bash
docker compose exec gateway-go sh
```
```bash
docker compose exec gateway-go bash
```


8. Stop containers without removing them:

```bash
docker compose stop
```

9. Start stopped containers again:

```bash
docker compose start
```

10. Remove unused Docker data (optional cleanup):

```bash
docker system prune
```

Use this carefully because it removes unused resources across Docker on your machine.

## Compose Health and Startup Checks
1. Gateway reports HTTP health after its listener is available.
2. Catalog reports both `/healthz` and `/readyz`; Compose watches the ready endpoint.
3. Order and user containers use lightweight process health checks.
4. Compose restart policies restart stopped containers unless explicitly stopped.
5. The gateway waits for order startup and catalog readiness before dependency startup.
6. The order service waits for user startup and catalog readiness before dependency startup.
7. Use `docker compose ps` to confirm health status rather than only listening ports.

## Container Troubleshooting
1. If the gateway returns 503, verify catalog and order container health in Compose.
2. If catalog readiness fails, inspect `docker compose logs catalog-service-py`.
3. If order creation fails, inspect catalog and user logs for the dependent gRPC call.
4. Invalid JSON, unknown fields, oversized bodies, and invalid identifiers return HTTP 400.
5. Rebuild the affected image after source or Dockerfile changes.
6. Remove stopped containers with `docker compose down` before a clean rebuild.
7. Persisted SQLite data is container-local unless explicitly mounted for a named volume.
