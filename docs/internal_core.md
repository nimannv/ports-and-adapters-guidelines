# Core Layer
**Path:** `internal/core/`

## Purpose
The `core` layer encapsulates the heart of the application's business logic, completely decoupled from external technologies (such as web frameworks, database drivers, or external APIs). It defines *what* the application does from a business perspective.

## Contains
- **[domain/](domain.md)**: Pure business entities, value objects, domain errors, and domain services.
- **[application/](internal_core_application.md)**: Orchestration logic and system use cases.
    - **[ports/](port.md)**: Interfaces defining contracts for inbound and outbound communication.
    - **[usecases/](usecase.md)**: Implementations of the inbound ports.

## Principles & Rules
- **Technology Agnostic**: This layer must not have any direct dependencies on external frameworks or infrastructure (e.g., `webservers`, `databases`, `message brokers`, ...). Only the language standard library is allowed.
- **High Cohesion, Low Coupling**: Code within the core is tightly related to the business problem, interacting with external layers strictly through defined ports.
- **Testable in Isolation**: The entire `core` layer must be unit-testable without requiring external infrastructure, using fakes or mocks for outbound ports.
- **Never Split by Technology**: The core is split by domain and by concept, never by technology (no `core/postgres/`, `core/http/`).

## Scaling
Split the core as a whole when there are separate business capabilities with their own language, data, and owners. It grows in four stages. Inside each stage, every layer follows its own rules: [domain/](domain.md#scaling), [usecases/](usecase.md#scaling), [ports/inbound/](inbound_port.md#scaling), [ports/outbound/](outbound_port.md#scaling).

### Stage 1: Flat
One small domain. Each layer is a flat list of files.

### Stage 2: Split Inside Layers
Several related concepts in one domain. Keep the layers on top and group by concept inside each layer.

```
internal/core/
├── domain/
│   ├── order/
│   └── customer/
└── application/
    ├── ports/
    │   ├── inbound/
    │   │   ├── order/
    │   │   └── customer/
    │   └── outbound/
    └── usecases/
        ├── order/
        └── customer/
```

### Stage 3: Split by Module
Separate bounded contexts. Put the **domain on top**, and give each module its own full set of layers.

```
internal/
├── core/
│   ├── order/                      # bounded context
│   │   ├── domain/
│   │   └── application/
│   │       ├── ports/
│   │       │   ├── inbound/
│   │       │   └── outbound/
│   │       └── usecases/
│   ├── billing/                    # bounded context
│   │   ├── domain/
│   │   └── application/
│   └── shared/                     # shared kernel (keep it tiny)
└── adapters/
    ├── inbound/
    │   ├── http/
    │   │   ├── public/             # scope
    │   │   │   ├── order/          # domain
    │   │   │   └── billing/
    │   │   └── admin/
    │   └── amqp/
    │       └── billing/
    └── outbound/
        ├── postgres/
        │   ├── order/
        │   └── billing/
        └── stripe/
```

- Modules talk to each other only through ports and never import each other's `domain/`.
- Keep shared code in a small `shared/` kernel. If it keeps growing, it is probably a missing module.

### Stage 4: Split into Services
A module needs to be deployed, scaled, or owned on its own. Since it already has clear ports and no shared domain, you can move it to its own repository or service with little rework.

See [Scale](../README.md#scale) for the general principles and [Enforcing Boundaries](internal.md#enforcing-boundaries) for keeping modules apart in Python.
