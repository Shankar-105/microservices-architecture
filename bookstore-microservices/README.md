# Bookstore Microservices (Polyglot)

This repository contains a phase-wise implementation of a polyglot microservices system for a bookstore.

## Current status
- Phase 1 completed: contracts and service scaffolding
- Services start with health endpoints and graceful shutdown skeletons

## Services
- `gateway-go` (Go): HTTP entrypoint
- `user-service-go` (Go): gRPC user service skeleton
- `order-service-go` (Go): gRPC order orchestration skeleton
- `catalog-service-py` (Python): gRPC + HTTP health skeleton
- `payment-service-py` (Python): gRPC + HTTP health skeleton

## Quick start
1. Review `docs/phase-01-foundation-and-contracts.md`
2. Generate protobuf stubs (requires `protoc`):
   - `make proto`
3. Run all services locally (no containers):
   - `make run-gateway`
   - `make run-user`
   - `make run-order`
   - `make run-catalog`
   - `make run-payment`

## Notes
- This phase intentionally keeps business logic minimal.
- Contract-first API design is locked in `proto/bookstore/`.
