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

| Step | Who | Maps |
|---|---|---|
| 1 | [Inbound adapter](inbound_adapter.md) | external format (JSON, protobuf, ...) → inbound adapter DTO → inbound port DTO |
| 2 | [Use case](usecase.md) | inbound port DTO → domain entities / value objects |
| 3 | [Use case](usecase.md) | domain entities / value objects → outbound port DTO |
| 4 | [Outbound adapter](outbound_adapter.md) | outbound port DTO → external format (ORM model, SQL row, API message, ...) |

Results travel back the same way, in reverse. Errors are translated the same way: outbound adapters turn technology errors into outbound port errors, use cases turn domain errors into inbound port errors, and inbound adapters turn inbound port errors into protocol responses.

## Scaling
- **[usecases/](usecase.md#scaling)**: by feature, then actor.
- **[ports/inbound/](inbound_port.md#scaling)**: follows the use-case split.
- **[ports/outbound/](outbound_port.md#scaling)**: by capability, never technology.

When the whole application layer needs to be split per bounded context, see [Core Scaling](internal_core.md#scaling).
