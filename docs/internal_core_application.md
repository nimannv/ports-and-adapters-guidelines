# Application Layer
**Path:** `internal/core/application/`

## Purpose
The `application` layer defines the application's specific use cases and orchestrates the flow of data and execution. It acts as an intermediary between the inbound adapters (external world) and the core domain logic.

## Contains
- **[ports/](port.md)**: Interfaces defining contracts for inbound and outbound communication.
    - **[inbound/](inbound_port.md)**: Interfaces defining what the system can do.
    - **[outbound/](outbound_port.md)**: Interfaces defining what the system needs from external systems.
- **[usecases/](usecase.md)**: Implementations of the inbound ports.

## Principles & Rules
- **Orchestration Layer**: Coordinates domain logic and interactions with external systems via ports.
- **Depends on Outbound Ports**: Depends on interfaces defined in `application/ports/outbound` to interact with external systems.
- **Depends on Domain**: Directly depends on the `domain` layer.
- **Translates at the Boundary**: Maps inbound port DTOs to domain objects, and domain objects to outbound port DTOs, so domain types never leak to adapters.
- **No Infrastructure Detail**: Remains independent of specific technologies (e.g., HTTP, gRPC, SQL).

## Data Flow
```
inbound adapter DTO ─► inbound port DTO ─► use case ─► domain entities / value objects
                                              │
                                              └─► outbound port DTO ─► outbound adapter DTO (ORM model, API message, ...)
```

## Scaling
Split `usecases/` by feature, then by actor, and let `ports/inbound/` follow the same split. Split `ports/outbound/` by the capability the core needs, never by technology. See the [Per-Layer Guide](../README.md#per-layer-guide).
