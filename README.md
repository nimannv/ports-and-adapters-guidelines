# Ports and Adapters Guidelines

A practical guide to implementing the Ports and Adapters (Hexagonal Architecture) pattern in Go and Python.

It decouples core business logic from external infrastructure, delivering systems that are testable, maintainable, and resilient to change.

In this approach, the focus is placed on strict decoupling and clarity of logic, never sacrificing clean boundaries for the sake of simpler code generation. With code generation becoming increasingly easier using agents, reducing boilerplate is no longer a reason to compromise architectural design. The primary goal is **predictability of implementation in each layer**—ensuring consistent, well-defined responsibilities and seamless maintainability.

![image](./diagram.svg)


## Folders Structure
- **[bootstrap/](docs/bootstrap.md)**: Application entry points (composition roots).
- **[configs/](docs/configs.md)**: Configuration templates and environment settings.
- **[scripts/](docs/configs.md#scripts)**: Utility and deployment scripts.
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

## Dependency Rule

Dependencies always point **inward**, toward the core:

- **Domain** depends on nothing.
- **Ports** depend on nothing.
- **Use cases** depend on the domain and on ports.
- **Adapters** depend on ports only.
- **Bootstrap** depends on everything, but only to wire it together.

The core never imports adapters or bootstrap code. See [internal/](docs/internal.md#dependency-rule) for the full table and how to enforce it in Go and Python.

## Data Flow

Each boundary has its own data types, and every layer maps them at its edge:

```
external format ─► inbound adapter DTO ─► inbound port DTO ─► use case ─► domain (entities, value objects)
                                                                 │
external format ◄─ outbound adapter DTO ◄─ outbound port DTO ◄──┘
```

Domain types never leave the core. See the [application layer](docs/internal_core_application.md#data-flow) for details.

## Domain Development

The architecture does not mandate a single domain modeling style. Depending on the complexity of the business logic, we have the flexibility to choose the most suitable approach.

### Examples:

- **Anemic / Data-Centric (Low Complexity)**: Simple data structures with minimal behavior, ideal for CRUD-heavy operations or simple workflows where logic lives primarily within use cases.
- **Rich Domain Model (Moderate Complexity)**: Entities encapsulate both **data and behavior**, enforcing their own invariants, validations, and state mutations.
- **Domain-Driven Design (High Complexity)**: Full tactical DDD (Aggregates, Entities, Value Objects, Domain Events, and Domain Services) to tackle intricate rules and strict consistency boundaries.
- **State Machine / Workflow-Driven**: Explicit state transitions and event-driven modeling for systems with complex lifecycles and strict state validation.
- ...

See [domain/](docs/domain.md) for the rules that apply whatever style you choose.

## Scale

The structure should stay small for a small service and grow with the codebase. Start flat and add another level only when one directory holds too many unrelated things. Each directory, and each layer inside `internal/`, can grow on its own.

The layout is the same for **Go and Python**. In the docs, `*` stands for the language extension (`main.*` means `main.go` or `main.py`). Use `snake_case` directory names (e.g. `place_order/`), because they are valid package names in both languages.

### Scaling Axes

Any directory can be split along three axes. Use them in this order:

1. **Domain**: business capability or bounded context (`order/`, `billing/`, `user/`).
2. **Scope**: audience, runtime, or deployable unit (`public/`, `admin/`, `internal/`, `api/`, `worker/`).
3. **Technology**: protocol or infrastructure (`http/`, `grpc/`, `postgres/`, `redis/`). Only adapters are split this way. Never split the core by technology.

### When to Split

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

### How Each Directory Scales

Each directory has its own split rules, described in its doc:

| Directory | Split by | Details |
|---|---|---|
| **bootstrap/** | Scope (runtime), then domain | [bootstrap.md](docs/bootstrap.md#scaling) |
| **configs/** | Environment, then component | [configs.md](docs/configs.md#scaling-configs) |
| **scripts/** | Purpose | [configs.md](docs/configs.md#scaling-scripts) |
| **core/** (as a whole) | Domain (bounded context), in growth stages | [internal_core.md](docs/internal_core.md#scaling) |
| **domain/** | Aggregate / business concept | [domain.md](docs/domain.md#scaling) |
| **usecases/** | Feature, then actor (or command/query) | [usecase.md](docs/usecase.md#scaling) |
| **ports/inbound/** | Same as use cases | [inbound_port.md](docs/inbound_port.md#scaling) |
| **ports/outbound/** | Capability the core needs, **never technology** | [outbound_port.md](docs/outbound_port.md#scaling) |
| **adapters/inbound/** | Technology → scope → domain (→ version) | [inbound_adapter.md](docs/inbound_adapter.md#scaling) |
| **adapters/outbound/** | Technology / external system → domain | [outbound_adapter.md](docs/outbound_adapter.md#scaling) |

### Other Directories

The same rules apply to every other directory as it appears, for example:

- **deployments/**: Split by target (`docker/`, `k8s/`, `terraform/`), then by environment.
- **api/**: Contract definitions (OpenAPI, Protobuf, AsyncAPI), split by protocol and then by version (`v1/`, `v2/`).
- **migrations/**: Split by datastore (`postgres/`, `mongo/`) when there is more than one.
- **test/**: Split by test type (`integration/`, `e2e/`, `fixtures/`). Unit tests stay next to the code they test.
- **docs/**: Split by audience or topic (`architecture/`, `adr/`, `runbooks/`).
- ...

> **Rule of thumb:** Don't create a level before you need it. Add one when a directory has too many unrelated files, or when a team or module needs a clear boundary.
