# Application Layer
**Path:** `internal/core/application/`

## Purpose
The `application` layer defines the application's specific use cases and orchestrates the flow of data and execution. It acts as an intermediary between the inbound adapters (external world) and the core domain logic.

## Contains
- **[ports/](port.md)**: Interfaces defining contracts for inbound and outbound communication.
- **[usecases/](usecase.md)**: Implementations of the inbound ports.

## Principles & Rules
- **Orchestration Layer**: Coordinates domain logic and interactions with external systems via ports.
- **Depends on Outbound Ports**: Depends on interfaces defined in `application/ports/outbound` to interact with external systems.
- **Depends on Domain**: Directly depends on the `domain` layer.
- **No Infrastructure Detail**: Remains independent of specific technologies (e.g., HTTP, gRPC, SQL).
