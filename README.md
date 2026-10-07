# Analytics Platform

Analytics Platform is an API-first product analytics service for collecting application events and querying event counts. It is designed for teams that need a small, project-scoped alternative to a full product analytics suite.

> **Project status:** early development. The repository currently contains the API gateway skeleton and product requirements. Event ingestion, project and API-key management, analytics queries, persistence, and API documentation are not implemented yet.

## Product workflow

The planned v0.1 workflow is:

1. Create a project to keep its data separate from other projects.
2. Create an API key for an operator or application.
3. Send application events to that project.
4. Query event counts for a selected time period.

The initial analytics scope is intentionally focused on counting events by name and time period. The product requirements are in [`docs/architecture/product.md`](docs/architecture/product.md).

## Technology

- Go 1.27.0
- Gin HTTP framework

## Prerequisites

- Go 1.27.0 or a compatible newer Go toolchain

## Run locally

From the repository root, start the API gateway with:

```sh
go run ./services/api-gateway/cmd/api-gateway
```

The current skeleton starts an HTTP server on `http://localhost:8080`. No application endpoints are available yet.

## Repository layout

```text
docs/architecture/       Product requirements and architecture notes
services/api-gateway/    API gateway service
```

## Development

Format Go files with `gofmt` before committing. Keep product behavior and API contracts aligned with the requirements in `docs/architecture/product.md`.

## License

No license has been specified yet.
