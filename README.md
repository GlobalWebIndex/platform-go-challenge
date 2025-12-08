# Assets Favourites API

A REST API service for managing user favourites of assets (charts, insights, audiences), implemented as part of the GWI Engineering Challenge.

## Overview

This service provides endpoints for users to manage their personal list of favourite assets. Assets can be of different types (charts, insights, audiences) with varying structures, and are created inline when adding favourites.

## Architecture

The service follows a domain-driven design with a layered, hexagonal-style architecture:

- **Domain layer** (`internal/domain`) - Pure business concepts independent of transport and storage
- **API layer** (`internal/api/v1`) - HTTP handlers, DTOs, and middleware
- **Repository pattern** - Storage abstractions allowing easy swapping of implementations (in-memory → database)

The architecture emphasizes separation of concerns, testability, and extensibility. Observability (logging, metrics, health checks) and authentication are treated as first-class concerns. See [04-architecture-and-standards.md](docs/04-architecture-and-standards.md) for detailed architecture and coding standards.

## Documentation

Detailed documentation is available in the `docs/` directory. Instead of full PRD (Product Requirements Document) and Technical Design Document templates, a lighter free format is used to keep things simple and reviewable without the overhead of formal templates.

- [00-challenge.md](docs/00-challenge.md) - Original GWI Engineering Challenge description
- [01-context.md](docs/01-context.md) - Project context, constraints, and non-functional requirements
- [02-domain-and-use-cases.md](docs/02-domain-and-use-cases.md) - Domain model and use cases
- [03-api-shape.md](docs/03-api-shape.md) - API design and specifications
- [04-architecture-and-standards.md](docs/04-architecture-and-standards.md) - Architecture and coding standards

## API Specification

The OpenAPI specification is available at [api/openapi.yaml](api/openapi.yaml).

## Status

This project is in the design phase. Implementation is pending.
