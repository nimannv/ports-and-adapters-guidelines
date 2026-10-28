# Adapters Layer
**Path:** `internal/adapters/`

## Purpose
The `adapters` layer houses all components responsible for interacting with the external world. It acts as the boundary between the application's core logic and external technologies, frameworks, or systems. Adapters translate external data formats and protocols into the language of the application's ports, and vice-versa.

## Contains
- **[inbound/](inbound_adapter.md)**: Driving adapters (e.g., HTTP, gRPC).
- **[outbound/](outbound_adapter.md)**: Driven adapters (e.g., SQL, MongoDB, External APIs).

## Principles & Rules
- **Outbound Adapters** Outbound adapters must implement the outbound ports defined in `internal/core/application/ports/outbound`.
- **Outbound Adapters DTOs** Outbound adapters use outbound ports DTOs as input and output. they map outbound ports DTOs to external data formats (e.g., SQL rows, MongoDB documents, external API responses) and vice versa.

- **Inbound Adapters** Inbound adapters use inbound ports defined in `internal/core/application/ports/inbound` to call usecases (inbound ports implementations).
- **Inbound Adapters DTOs** Inbound adapters use their own DTOs to interact with external systems (e.g., HTTP JSON, gRPC messages) and map them to inbound ports DTOs.

- **Technology Specific**: This is the only layer where technology-specific code (e.g., `net/http` package, gRPC proto-generated code) is allowed and expected.
- **No Core Logic**: Adapters contain minimal to no business logic. Their primary role is translation and interaction with external systems.
