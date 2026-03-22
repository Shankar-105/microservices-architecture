# Explaining the Zrchitecture and Microservices

I want to be very clear from the start: this repository is not production ready.
I built this as a learning-first microservices project to explain how microservices are designed, how service boundaries are chosen, and how services communicate reliably.

## Why I Built It This Way
When people learn microservices, they often jump directly into tools and deployments, but skip architecture decisions.
In this repo, I wanted to keep the flow small and understandable:
1. One gateway receives HTTP requests.
2. Internal services communicate via gRPC.
3. Each service owns its own responsibility.

That gives us a clean path to understand the core idea: microservices are not about having many folders, they are about clear boundaries and independent evolution.

## Services In This Repo
1. Gateway service in Go.
2. Order service in Go.
3. User service in Go
4. Catalog service in Python.

## Why Go For Gateway and Order
I used Go for the gateway and order service because Go is excellent for backend services that need high concurrency and low overhead.
For example, a gateway can receive many simultaneous user requests and Go's goroutine model handles this naturally with very low memory cost.
That makes Go a strong fit for request routing, protocol translation, and orchestration-heavy workloads.

Another reason is predictability in production-like server behavior:
1. Fast startup.
2. Strong static typing.
3. Easy deployment as single binaries.

## Why Python For Catalog
I used Python for the catalog side because recommendation-heavy domains often evolve toward data science and machine learning workflows.
Even if this repo currently keeps catalog logic simple with just returning the books, Python gives a natural path if later we add:
1. Recommendation models.
2. Feature engineering pipelines.
3. ML-serving integrations.

This demonstrates a practical microservices principle: pick the right tool per service domain, not one language for everything by force or to showcase that i have used many languages.

## When To Use Microservices
Microservices are useful when the **Two main reasons**:

1. Teams need independent deployment of parts of the system.
2. Different technical stacks make sense for different domains.

Example:
If your order workflow changes weekly but your user profile service is stable, microservices allow changing order-service independently.

## When Not To Use Microservices
Do not start with microservices when:
1. Product scope is still small and changing rapidly, that is your just starting off with your project.
2. Team is very small and operational overhead will hurt speed.
3. You do not yet need independent scaling/deployments.

In those cases, a modular monolith is usually better, simpler, and cheaper to run and alwasy the goto choice.

## Why gRPC In This Project
I used gRPC for internal service-to-service communication because it is contract-first and language-agnostic.

Benefits in this repo:
1. Shared `.proto` contracts reduce integration ambiguity.
2. Generated code keeps Go and Python in sync.
3. Binary transport is efficient for backend-to-backend calls.
4. Schema changes can be managed more safely than ad-hoc JSON.

This is a good use of gRPC: internal calls where performance, contracts, and multi-language support matter.

## Final Thought
This repository is intentionally small, but rich in architecture learning.
The goal is to understand why microservices are designed in a certain way and when to use them.