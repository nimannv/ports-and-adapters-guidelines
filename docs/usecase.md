# Use Cases
**Path:** `internal/core/application/usecases/`

## Purpose
The `usecases` package contains concrete implementations of the [inbound ports](./inbound_port.md). It defines *how* the application performs its specific tasks (application layer logic) by orchestrating domain logic, using outbound ports and other usecases.

## Contains
- **Usecase Implementations**: One file per use case (e.g., `place_order.*`), implementing an interface defined in [inbound ports](./inbound_port.md).

## Principles & Rules

- **Dependency Direction**:
    - **Depends on Domain**: Directly uses entities, value objects, and domain services.
    - **Depends on Outbound Ports**: Depends on interfaces defined in `application/ports/outbound` to interact with external systems.
    - **Depends on other Use Cases**: Depends on interfaces defined in `application/ports/inbound` to use other usecases (never on their concrete implementations).

- **Orchestration Layer**:
    - Coordinates domain logic
    - Interactions with external systems via outbound ports
    - Interactions with other usecases

- **DTO and Error Mapping**:
    - Maps inbound port DTOs to domain objects, and domain objects to outbound port DTOs (and back).
    - Translates domain errors into the errors declared by the inbound port.

- **Transactional Consistency**:
    - Often serves as the boundary for database transactions. Express the transaction through an outbound port (e.g., a `UnitOfWork` / `TxManager`), never through a database driver.

- **Dependency Injection**:
    - Receives its dependencies (outbound ports, other use cases) through its constructor. They are wired in [bootstrap](bootstrap.md).

- **Testable in Isolation**:
    - Unit-tested with in-memory fakes or mocks of the outbound ports. No database, no network.

## Scaling

### When to Split
- **Too many files**: the flat directory is hard to scan (roughly 10–15 use cases).
- **Different features**: use cases group naturally around a business concept (orders, customers, payments).
- **Different actors**: some use cases serve customers, others admins or internal jobs, with different permissions and rules.
- **Reads and writes diverge**: queries need different models or data sources than commands.

### How to Split
Split in this order, and stop as soon as the directory is readable again:

1. **By feature**: group by business concept, matching [domain/](domain.md#scaling) (`usecases/order/`, `usecases/customer/`).
2. **By actor**: inside a feature, or at the top when actors are very different (`usecases/admin/`, `usecases/customer/`).
3. **By command/query**: only if reads and writes differ a lot (`usecases/order/commands/`, `usecases/order/queries/`).

```
# 1. Flat                  # 2. By feature              # 3. By feature, then actor
usecases/                  usecases/                    usecases/
├── place_order.*   ──►    ├── order/            ──►    ├── order/
├── cancel_order.*         │   ├── place_order.*        │   ├── place_order.*
├── refund_order.*         │   ├── cancel_order.*       │   ├── cancel_order.*
└── register_customer.*    │   └── refund_order.*       │   └── admin/
                           └── customer/                │       └── refund_order.*
                               └── register_customer.*  └── customer/
                                                            └── register_customer.*
```

- [Inbound ports](inbound_port.md#scaling) always mirror this tree.
- Name files after the action (`place_order.*`), not the entity (`order_service.*`). One use case per file.
- When features turn into separate bounded contexts, split the whole core instead (see [Core Scaling, Stage 3](internal_core.md#stage-3-split-by-module)).

See [Scale](../README.md#scale) for the general principles.
