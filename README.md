# Ports and Adapters Guidelines

A practical guide to implementing the Ports and Adapters (Hexagonal Architecture) pattern in Go and Python.

It decouples core business logic from external infrastructure, delivering systems that are testable, maintainable, and resilient to change.


In this approach, the focus is placed on strict decoupling and clarity of logic, never sacrificing clean boundaries for the sake of simpler code generation. With code generation becoming increasingly easier using agents, reducing boilerplate is no longer a reason to compromise architectural design. The primary goal is **predictability of implementation in each layer**—ensuring consistent, well-defined responsibilities and seamless maintainability.

![image](./diagram.svg)


## Folders Structure
- **[bootstrap/](docs/bootstrap.md)**: Application entry points (composition roots).
        - ...
- **[configs/](docs/configs.md)**: Configuration templates and environment settings.
- **[scripts/](docs/configs.md)**: Utility and deployment scripts.
- **[internal/](docs/internal.md)**: Main application code, including core logic and adapters.
    - **[core/](docs/internal_core.md)**: The heart of the application, containing domain logic and application use cases.
        - **[domain/](docs/domain.md)**: Pure business entities, value objects, and domain services.
        - **[application/](docs/internal_core_application.md)**: Orchestration logic and system use cases.
            - **[ports/](docs/port.md)**: Interfaces defining contracts for inbound and outbound communication.
                - **[inbound/](docs/inbound_port.md)**: Interfaces defining what the system can do.
                - **[outbound/](docs/outbound_port.md)**: Interfaces defining what the system needs from external systems.
            - **[usecases/](docs/usecase.md)**: Implementations of the inbound ports.
    - **[adapters/](docs/adapters.md)**: Interaction with the external world.
        - **[inbound/](docs/inbound_adapter.md)**: Driving adapters (e.g., HTTP, gRPC, CLI, AMQP, ...).
        - **[outbound/](docs/outbound_adapter.md)**: Driven adapters (e.g., SQL, MongoDB, Redis, External APIs, AMQP, ...).

## Domain Development

The architecture does not mandate a single domain modeling style. Depending on the complexity of the business logic, we have the flexibility to choose the most suitable approach.

### Examples:

- **Anemic / Data-Centric (Low Complexity)**: Simple data structures with minimal behavior, ideal for CRUD-heavy operations or simple workflows where logic lives primarily within use cases.
- **Rich Domain Model (Moderate Complexity)**: Entities encapsulate both **data and behavior**, enforcing their own invariants, validations, and state mutations.
- **Domain-Driven Design (High Complexity)**: Full tactical DDD (Aggregates, Entities, Value Objects, Domain Events, and Domain Services) to tackle intricate rules and strict consistency boundaries.
- **State Machine / Workflow-Driven**: Explicit state transitions and event-driven modeling for systems with complex lifecycles and strict state validation.
- ...

## Scale

The structure should stay small for a small service and grow with the codebase. Start flat and add another level only when one directory holds too many unrelated things. Each directory, and each layer inside `internal/`, can grow on its own.

The layout is the same for **Go and Python**. In the trees below, `*` stands for the language extension (`main.*` means `main.go` or `main.py`). Use `snake_case` directory names (e.g. `place_order/`), because they are valid package names in both languages.

### Scaling Axes

Any directory can be split along three axes. Use them in this order:

1. **Domain**: business capability or bounded context (`order/`, `billing/`, `user/`).
2. **Scope**: audience, runtime, or deployable unit (`public/`, `admin/`, `internal/`, `api/`, `worker/`).
3. **Technology**: protocol or infrastructure (`http/`, `grpc/`, `postgres/`, `redis/`). Only adapters are split this way. Never split the core by technology.

### bootstrap/

Start with one entry point. Each entry point stays a thin composition root: it wires adapters to use cases and starts the program. As the codebase grows, split it by **scope** (runtime), then by **domain**.

```
# 1. Single entry point     # 2. Split by scope (runtime)    # 3. Split by domain, then scope
bootstrap/                  bootstrap/                       bootstrap/
└── main.*                  ├── api/                         ├── order/
                            │   └── main.*                   │   ├── api/
                            ├── worker/                      │   │   └── main.*
                            │   └── main.*                   │   └── worker/
                            ├── cli/                         │       └── main.*
                            │   └── main.*                   ├── billing/
                            └── migrate/                     │   ├── api/
                                └── main.*                   │   └── worker/
                                                             └── admin/
                                                                 └── cli/
                                                                     └── main.*
```

- **Split by scope** when one program has to start different kinds of work (serve HTTP, consume a queue, run a job). These need to be deployed, scaled, or restarted separately.
- **Split by domain** when separate business modules get their own deployments, release cycles, or owning teams. Each one wires only the use cases and adapters it needs.

> **Naming:** Go projects usually call this directory `cmd/`. We use `bootstrap/` because the name works in both languages and says what the code does: it builds the application and starts it. In Python, `cmd` is a standard library module, and a top-level package named `cmd` would shadow it.

### configs/

Start with one config file. As the number of environments and components grows, split the files by **environment** and then by **component**. Defaults live in one base file, and each environment overrides only what it needs.

```
configs/                    configs/
└── config.yaml      ──►    ├── base.yaml             # shared defaults
                            ├── local/
                            │   ├── app.yaml
                            │   └── database.yaml
                            ├── staging/
                            │   └── app.yaml
                            └── production/
                                ├── app.yaml
                                └── database.yaml
```

### scripts/

Start with a few scripts at the top level. As they pile up, group them by **purpose**. Keep a single entry point (`Makefile` / `Taskfile`) so people don't need to know the internal layout.

```
scripts/                    scripts/
├── build.sh         ──►    ├── build/
├── test.sh                 │   ├── docker.sh
└── deploy.sh               │   └── package.sh
                            ├── ci/
                            │   ├── lint.sh
                            │   └── test.sh
                            ├── deploy/
                            │   ├── staging.sh
                            │   └── production.sh
                            └── dev/
                                ├── seed.sh
                                └── setup.sh
```

### internal/

The core and the adapters scale on their own, because the ports keep them apart. The real question is **when** to split each layer and **based on what**.

#### General Criteria

Split a directory when one or more of these signals appear:

| Signal | Question to ask |
|---|---|
| **Cohesion** | Do these files change together? If not, they don't belong together. |
| **Language** | Does the same word mean different things in different places (e.g. `Customer` in sales vs. billing)? That marks a domain boundary. |
| **Ownership** | Do different teams or people own different parts? |
| **Rate of change** | Does one part change every week while another stays stable? |
| **Non-functional needs** | Do parts need different auth, SLAs, scaling, or deployment? That marks a scope boundary. |
| **Cognitive load** | Can you still understand the directory at a glance? Roughly 10–15 files is a sign to group. |

Don't split:

- **Into technical buckets** (`entities/`, `enums/`, `value_objects/`, `utils/`). Group by business concept instead.
- **When it creates a dependency cycle.** A cycle means the boundary is wrong. Merge the parts or extract a shared concept.
- **Just in case.** An empty level is noise.

#### Per-Layer Guide

| Layer | Split when... | Split by | Example |
|---|---|---|---|
| **core/** (as a whole) | There are separate business capabilities with their own language, data, and owners. | Domain (bounded context) | `core/order/`, `core/billing/` |
| **domain/** | Entities form separate groups with their own rules and change for different reasons. | Aggregate / business concept | `domain/order/`, `domain/customer/` |
| **usecases/** | There are many use cases, or they serve different actors or features. | Feature, then actor (or command/query if reads and writes differ a lot) | `usecases/order/place_order.*`, `usecases/admin/...` |
| **ports/inbound/** | Use cases get split. Inbound ports follow the use-case split. | Same as use cases | `ports/inbound/order/` |
| **ports/outbound/** | An interface has consumers that each use only part of it (interface segregation). | Capability the core needs, **never technology** | `OrderRepository`, `PaymentGateway`, `Notifier` |
| **adapters/inbound/** | A protocol has many handlers, or audiences need different auth or middleware. | Technology → scope → domain (→ version) | `http/public/order/`, `http/admin/`, `grpc/v1/` |
| **adapters/outbound/** | One adapter implements ports for several domains, or talks to several external systems. | Technology / external system → domain | `postgres/order/`, `stripe/` |

#### Growth Stages of the Core

**Stage 1 (Flat):** One small domain. Each layer is a flat list of files.

**Stage 2 (Split inside layers):** Several related concepts in one domain. Keep the layers on top and group by concept inside each layer.

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

**Stage 3 (Split by module):** Separate bounded contexts. Put the **domain on top**, and give each module its own full set of layers. Modules talk to each other only through ports and never import each other's `domain/`. Keep shared code in a small `shared/` kernel.

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

**Stage 4 (Split into services):** A module needs to be deployed, scaled, or owned on its own. Since it already has clear ports and no shared domain, you can move it to its own repository or service with little rework.

> **Python note:** Go's compiler enforces `internal/`. Python doesn't, so enforce these boundaries with a tool such as [import-linter](https://github.com/seddonym/import-linter) (e.g. a contract that says `core` must not import `adapters`, and `order.domain` must not import `billing.domain`).

### Other directories

The same rule applies to every other directory as it appears, for example:

- **deployments/**: Split by target (`docker/`, `k8s/`, `terraform/`), then by environment.
- **api/**: Contract definitions (OpenAPI, Protobuf, AsyncAPI), split by protocol and then by version (`v1/`, `v2/`).
- **migrations/**: Split by datastore (`postgres/`, `mongo/`) when there is more than one.
- **test/**: Split by test type (`integration/`, `e2e/`, `fixtures/`). Unit tests stay next to the code they test.
- **docs/**: Split by audience or topic (`architecture/`, `adr/`, `runbooks/`).
- ...

> **Rule of thumb:** Don't create a level before you need it. Add one when a directory has too many unrelated files, or when a team or module needs a clear boundary.

