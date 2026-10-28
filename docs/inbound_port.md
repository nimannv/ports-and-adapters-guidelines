# Inbound Ports
**Path:** `internal/core/application/ports/inbound/`

## Purpose
**Inbound ports** define the external interfaces for interacting with the application's core logic. They represent the system [usecases](usecase.md) and their DTOs. Such as `CreateUser`, `UpdateUser`, `DeleteUser`, etc.

## Contains
- **Interfaces**: Definitions of usecase interfaces that can be called by [inbound adapters](./inbound_adapter.md) or other usecases.
- **execute method**: Each usecase interface has only one public method called `exec`. It takes input DTO as parameter and returns output DTO or error.
- **DTOs**: Data Transfer Objects used to pass input/output data between [inbound adapters](./inbound_adapter.md) and the core. 
    > DTOs are the main usage of inobound port, otherwise it would be only an interface with `exec` method.

## Principles & Rules
- **Contract Definition**: Define what the system can do and how other parts can use it.
- **Implementation**: These interfaces are implemented by [usecases](usecase.md).
- **Dependency Direction**: [Inbound adapters](inbound_adapter.md) (e.g., HTTP handlers) depend on and call these interfaces.
- **Technology Independent**: They must not reference any specific technology (e.g., HTTP status codes or gRPC metadata) this is [inbound adapters](inbound_adapter.md) responsibility.
