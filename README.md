# Bookstore Microservices Learning Repo

## Important Note
This repository is **not a production-ready microservices system**.
It is a practical learning example to understand how microservices work, how services communicate, and how to design boundaries with different languages.

## What This Repo Demonstrates
This project shows a **polyglot microservices approach**:
1. A **Go gateway service** for HTTP entry and request routing.
2. A **Go order service** for order orchestration and persistence.
3. A **Python catalog service** for catalog/business logic.
4. **gRPC + Protocol Buffers** for internal service-to-service communication.

## Why Polyglot Here
I intentionally used both Go and Python to show when teams may pick different languages per service responsibility.
1. Go is a strong fit for gateway-style high-concurrency request handling as goroutines use all the cpu cores without any hesitation.
2. Python is a strong fit for data/ML-friendly domains where recommendation logic can evolve faster so i used it in the catalog-service but for now i didn't use any ML model for books recommendations.
3. In real systems, language choice should follow team strengths, runtime requirements, and long-term maintainability.

## Why gRPC Here
gRPC is used for internal service communication because it provides:
1. Strong contracts via `.proto` files.
2. Type-safe generated clients/servers in multiple languages.
3. Efficient binary transport for service-to-service calls.

## Read This Next
If you want to know about the architecture and what tech stack is used where and why:
1. [docs/architecture.md](docs/architecture.md)

If you want to understand actual runtime flow of this repo, including all 4 routes:
1. [docs/repo-working.md](docs/repo-working.md)

If you want setup using Makefile:
1. [docs/setup-with-makefile.md](docs/setup-with-makefile.md)

If you want setup using Docker:
1. [docs/setup-with-docker.md](docs/setup-with-docker.md)


# Thank you 
**_M.Bhavani Shankar Student at Anits Collage Vizag Andhra Pradesh India_**