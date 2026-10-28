# Internal Packages
**Path:** `internal/`

## Purpose
The `internal` package contains the core business logic and all implementation details of the application. Everything outside it ([bootstrap/](bootstrap.md), [configs/ and scripts/](configs.md), ...) only builds, configures, or runs what lives here.

## Contains
- **[core/](internal_core.md)**: The heart of the application, containing domain logic and application use cases.
- **[adapters/](adapters.md)**: Interaction with the external world.
    - **[inbound/](inbound_adapter.md)**: Driving adapters (e.g., HTTP, gRPC, CLI, AMQP, ...).
    - **[outbound/](outbound_adapter.md)**: Driven adapters (e.g., SQL, MongoDB, Redis, External APIs, AMQP, ...).

## Principles & Rules
- **Encapsulation**: All application-specific logic must reside here.
- **Separation of Concerns**: Business logic (core) is kept strictly separate from technical details (adapters).
- **Dependency Rule**: Dependencies always point inward, toward the core (see below).

## Dependency Rule

| Package | May depend on |
|---|---|
| `core/domain` | Nothing (standard library only) |
| `core/application/ports/inbound` | Nothing (standard library only) |
| `core/application/ports/outbound` | Nothing (standard library only) |
| `core/application/usecases` | `domain`, `ports/inbound`, `ports/outbound` |
| `adapters/inbound` | `ports/inbound` |
| `adapters/outbound` | `ports/outbound` |
| `bootstrap/` (outside `internal/`) | Everything, for wiring only |

- The core never imports adapters or bootstrap code.
- Adapters never import `domain` or `usecases`.
- Adapters never import each other; an inbound adapter reaches an outbound adapter only through a use case.

## Enforcing Boundaries
- **Go**: The compiler only stops code outside the module from importing `internal/`. It does not check the layer rules inside it. Use a linter such as [depguard](https://github.com/OpenPeeDeeP/depguard) (via `golangci-lint`) or [go-arch-lint](https://github.com/fe3dback/go-arch-lint).
- **Python**: The language enforces nothing. Make every directory a regular package (with an `__init__.py`) and check the rules with [import-linter](https://github.com/seddonym/import-linter), e.g. contracts that say `internal.core` must not import `internal.adapters`, `internal.adapters` must not import `internal.core.domain`, and `order.domain` must not import `billing.domain`.

## Scaling
The core and the adapters scale on their own, because the ports keep them apart:

- **[core/](internal_core.md#scaling)**: grows in stages, from flat layers to bounded-context modules.
- **[adapters/](adapters.md#scaling)**: split by technology first, then scope and domain.

See [Scale](../README.md#scale) for the general principles.
