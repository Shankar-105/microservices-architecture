# Repo Working

This document explains how to call each route and what happens internally in this repository.

## Routes In This Repo
1. `GET /healthz`
2. `GET /getallbooks`
3. `POST /createorder`
4. `GET /getorders?user_id=...`

Gateway base URL for local run: `http://localhost:8080`

## 1) GET /healthz

### How To Call
1. Browser:
	- Open `http://localhost:8080/healthz`
2. curl:
	- `curl http://localhost:8080/healthz`
3. Postman:
	- Method: `GET`
	- URL: `http://localhost:8080/healthz`

### What Happens Internally
1. Request reaches the Gateway service.
2. Gateway handles it directly (no gRPC call needed).
3. Gateway returns simple JSON health response.

## 2) GET /getallbooks

### How To Call
1. Browser:
	- Open `http://localhost:8080/getallbooks`
2. curl:
	- `curl http://localhost:8080/getallbooks`
3. Postman:
	- Method: `GET`
	- URL: `http://localhost:8080/getallbooks`

### What Happens Internally
1. Request enters Gateway.
2. Gateway calls Catalog service over gRPC (`GetAllBooks`).
3. Catalog service reads books from its own SQLite DB.
4. Catalog service returns list of books to Gateway.
5. Gateway maps the response to HTTP JSON and sends it back.

## 3) POST /createorder

### How To Call
1. Browser:
	- Use Postman or curl for this route (because it is POST + JSON body).
2. curl:
	- `curl -X POST http://localhost:8080/createorder -H "Content-Type: application/json" -d '{"user_id":"u-100","book_id":"b-100","quantity":1}'`
3. Postman:
	- Method: `POST`
	- URL: `http://localhost:8080/createorder`
	- Body (JSON):
	  - `{"user_id":"u-100","book_id":"b-100","quantity":1}`

### What Happens Internally
1. Request enters Gateway and payload is validated.
2. Gateway calls Order service over gRPC (`CreateOrder`).
3. Order service checks user via User service gRPC.
4. Order service checks selected book via Catalog service gRPC.
5. If valid, Order service writes order to its own DB.
6. Order service returns order id, status, and total.
7. Gateway returns HTTP JSON response.

## 4) GET /getorders?user_id=...

### How To Call
1. Browser:
	- Open `http://localhost:8080/getorders?user_id=u-100`
2. curl:
	- `curl "http://localhost:8080/getorders?user_id=u-100"`
3. Postman:
	- Method: `GET`
	- URL: `http://localhost:8080/getorders?user_id=u-100`

### What Happens Internally
1. Request enters Gateway.
2. Gateway validates `user_id` query parameter.
3. Gateway calls Order service over gRPC (`GetOrders`).
4. Order service reads orders for that user from its DB.
5. Order service returns the order list.
6. Gateway converts result to HTTP JSON and returns it.

## Quick Data Ownership Summary
1. Gateway does not own business data; it routes and maps.
2. Order service owns order records.
3. Catalog service owns book records.
4. User service owns user validation data.

## Defensive Request Behavior
1. Gateway JSON bodies are limited to 1 MiB and unknown fields are rejected.
2. An order body must contain exactly one JSON object.
3. User and book identifiers are trimmed and limited to 128 characters.
4. Order quantities must be between 1 and 10000.
5. The gateway applies a five-second deadline to downstream gRPC calls.
6. Catalog and order lookups translate gRPC status codes into HTTP status codes.
7. NOT_FOUND responses become HTTP 404 and unavailable dependencies become HTTP 503.
8. Invalid arguments become HTTP 400 and failed preconditions become HTTP 409.
9. Every gateway response includes an `X-Request-ID` header for tracing.
10. The recovery middleware converts unexpected panics into JSON HTTP 500 responses.

## Failure Scenarios
1. An unavailable user service produces a dependency failure rather than a created order.
2. An unknown user or book is reported without writing an order.
3. An unavailable catalog book fails before repository insertion.
4. Integer overflow is checked before computing an order total.
5. Order history is capped at 500 records to limit response and memory growth.
6. Repository state and database contexts are checked before data operations.
7. Catalog repository exceptions are isolated from the network response boundary.
