# Use Cases
**Path:** `internal/core/application/usecases/`

## Purpose
The `usecases` package contains concrete implementations of the [inbound ports](./inbound_port.md). It defines *how* the application performs its specific tasks (application layer logic) by orchestrating domain logic, using outbound ports and other usecases.

## Contains
- **Usecase Implementations**: Go files implementing the interfaces defined in [inbound ports](./inbound_port.md).

## Principles & Rules

- **Dependency Direction**:
    - **Depends on Domain**: Directly uses entities, value objects, and domain services.
    - **Depends on Outbound Ports**: Depends on interfaces defined in `application/ports/outbound` to interact with external systems.
    - **Depends on other Use Cases**: Depends on interfaces defined in `application/ports/inbound` to use other usecases.

- **Orchestration Layer**:
    - Coordinates domain logic
    - Interactions with external systems via outbound ports
    - Interactions with other usecases

- **Transactional Consistency**:
        - Often serves as the boundary for database transactions.
