# Outbound Ports
**Path:** `internal/core/application/ports/outbound/`

## Purpose
Outbound ports define the external systems or capabilities that the application core needs from the outside world. This includes databases, message queues, external APIs, or file systems.

## Contains
- **Interfaces**: Definitions of methods that the core calls to interact with external systems.
- **DTOs**: Data Transfer Objects used to pass data back and forth between the core and external systems.

## Principles & Rules
- **Dependency Inversion**: The core defines its own needs as ports.
- **Technology Agnostic**: Ports define *what* is needed, not *how* it's done (e.g., `UserRepository` instead of `SQLUserRepository`).
- **Implementation**: These interfaces are implemented by [Outbound Adapters](outbound_adapter.md).
- **Error Handling**: Ports define error types or use generic domain-specific error types.
