# Bookstore Microservices (Polyglot)

This repository contains a polyglot microservices system for a bookstore.
## Services
- `gateway-go` (Go): HTTP entrypoint
- `user-service-go` (Go): gRPC user service skeleton
- `order-service-go` (Go): gRPC order orchestration skeleton
- `catalog-service-py` (Python): gRPC + HTTP health skeleton
- `payment-service-py` (Python): gRPC + HTTP health skeleton

## Quick start
1. Generate protobuf stubs (requires `protoc`):
   - `make proto`
2. Run all services locally (no containers):
   - `make run-all`
3. Stop all at once:
   - `make stop-all` 
## Notes
Every thing to be updated soon!