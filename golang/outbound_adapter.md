# Outbound Adapters
**Path:** `internal/adapters/outbound/`

## Purpose
**Outbound adapters** implement the [outbound ports](../outbound_port.md) defined by the application's core. They fulfill the core's needs by interacting with external systems, such as databases, APIs, or message queues, 

## Contains
- **repository**: Implementations of data persistence outbound ports (e.g., PostgreSQL, MongoDB).
- **external_api**: Implementations of outbound ports for calling third-party services.
- **any other outbound adapters**...

## Principles & Rules
- **Implements Outbound Ports**: Outbound adapters must implement the [outbound ports](../outbound_port.md) defined by the application's core.
- **DTO Mapping**: They use the outbound ports DTOs and map them to the external systems DTOs (HTTP messages, ORM models, etc.).
- **Technology Specific**: They are the bridge to specific external systems and contain technology-specific logic (e.g., SQL queries, API clients).
- **Error Translation**: Responsible for translating technology-specific errors (e.g., database connection timeouts) into domain-specific error types that the core can understand.
- **No Business Logic**: Their primary role is translation and interaction, not containing business rules.
