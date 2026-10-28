# Adapters Layer
**Path:** `internal/adapters/`

## Purpose
The `adapters` layer houses all components responsible for interacting with the external world. It acts as the boundary between the application's core logic and external technologies, frameworks, or systems. Adapters translate external data formats and protocols into the language of the application's ports, and vice-versa.

## Contains
- **[inbound/](inbound_adapter.md)**: Driving adapters (e.g., HTTP, gRPC, CLI, AMQP, ...).
- **[outbound/](outbound_adapter.md)**: Driven adapters (e.g., SQL, MongoDB, Redis, External APIs, AMQP, ...).

## Principles & Rules
- **Outbound Adapters**: Outbound adapters must implement the outbound ports defined in `internal/core/application/ports/outbound`.
- **Outbound Adapters DTOs**: Outbound adapters use outbound port DTOs as input and output. They map outbound port DTOs to external data formats (e.g., SQL rows, MongoDB documents, external API responses) and vice versa.

- **Inbound Adapters**: Inbound adapters use inbound ports defined in `internal/core/application/ports/inbound` to call usecases (inbound port implementations).
- **Inbound Adapters DTOs**: Inbound adapters use their own DTOs to interact with external systems (e.g., HTTP JSON, gRPC messages) and map them to inbound port DTOs.

- **Ports Only**: Adapters depend on ports only. They never import the `domain` or `usecases` packages.
- **No Adapter-to-Adapter Calls**: An inbound adapter never calls an outbound adapter directly; it always goes through a use case.
- **Technology Specific**: This is the only layer where technology-specific code (e.g., Go `net/http`, `database/sql`, gRPC generated code; Python FastAPI, SQLAlchemy, psycopg) is allowed and expected.
- **No Core Logic**: Adapters contain minimal to no business logic. Their primary role is translation and interaction with external systems.

## Scaling
Adapters are the only layer split by technology:

- **inbound/**: technology → scope → domain (→ version), e.g. `http/public/order/`, `http/admin/`, `grpc/v1/`.
- **outbound/**: technology / external system → domain, e.g. `postgres/order/`, `stripe/`.

See the [Per-Layer Guide](../README.md#per-layer-guide).
