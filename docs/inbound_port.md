# Inbound Ports
**Path:** `internal/core/application/ports/inbound/`

## Purpose
**Inbound ports** define the external interfaces for interacting with the application's core logic. They represent the system [usecases](usecase.md) and their DTOs. Such as `CreateUser`, `UpdateUser`, `DeleteUser`, etc.

Inbound ports exist for three reasons:

- **DTOs**: Each use case gets its own input and output types. They shape the data that crosses into the core and keep transport models (HTTP requests, gRPC messages, ...) and domain entities out of the contract.
- **A single list of what the core offers**: All use cases are listed in one place. To see everything the system can do, read this directory. You don't need to dig through handlers or use-case implementations.
- **One contract for many adapters**: HTTP, gRPC, CLI, AMQP and other [inbound adapters](./inbound_adapter.md) all call the same interface with the same DTOs. A new adapter reuses the existing contract, and the core doesn't change.

## Contains
- **Interfaces**: Definitions of usecase interfaces that can be called by [inbound adapters](./inbound_adapter.md) or other usecases.
- **execute method**: Each usecase interface has only one public method called `exec`. It takes input DTO as parameter and returns output DTO or error.
- **DTOs**: Data Transfer Objects used to pass input/output data between [inbound adapters](./inbound_adapter.md) and the core. 
    > DTOs are the main reason for inbound ports (see [Purpose](#purpose)). Without them, a port would be just an interface with an `exec` method.

## Principles & Rules
- **Contract Definition**: Define what the system can do and how other parts can use it.
- **Implementation**: These interfaces are implemented by [usecases](usecase.md).
- **Dependency Direction**: [Inbound adapters](inbound_adapter.md) (e.g., HTTP handlers) depend on and call these interfaces.
- **Technology Independent**: They must not reference any specific technology (e.g., HTTP status codes or gRPC metadata) this is [inbound adapters](inbound_adapter.md) responsibility.
