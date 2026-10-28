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
- **Dependency Rule**: Dependencies always point inward, toward the core. The core never imports adapters or bootstrap code.

| Package | May depend on |
|---|---|
| `core/domain` | Nothing (standard library only) |
| `core/application/ports/inbound` | Nothing (standard library only) |
| `core/application/ports/outbound` | Nothing (standard library only) |
| `core/application/usecases` | `domain`, `ports/inbound`, `ports/outbound` |
| `adapters/inbound` | `ports/inbound` |
| `adapters/outbound` | `ports/outbound` |
| `bootstrap/` (outside `internal/`) | Everything, for wiring only |

## Enforcing Boundaries
- **Go**: The compiler only stops code outside the module from importing `internal/`. It does not check the layer rules inside it. Use a linter such as [depguard](https://github.com/OpenPeeDeeP/depguard) (via `golangci-lint`) or [go-arch-lint](https://github.com/fe3dback/go-arch-lint).
- **Python**: Nothing is enforced by the language. Every directory is a regular package (with an `__init__.py`), and the rules are checked with [import-linter](https://github.com/seddonym/import-linter):

```ini
# .importlinter
[importlinter]
root_package = internal

[importlinter:contract:layers]
name = Domain is the innermost layer
type = layers
layers =
    internal.adapters
    internal.core.application
    internal.core.domain

[importlinter:contract:adapters-use-ports-only]
name = Adapters depend on ports only
type = forbidden
source_modules = internal.adapters
forbidden_modules =
    internal.core.domain
    internal.core.application.usecases
```

## Conventions Used in Examples
All docs share one running example: a small shop that lets a customer **place an order**. It uses the flat layout (Stage 1 in [Growth Stages of the Core](../README.md#growth-stages-of-the-core)).

- **Go**: module `example.com/shop`; files under `internal/...` end in `.go`.
- **Python**: packages are imported from the repository root (`internal.core.domain.order`); files end in `.py`.

| Layer | Go | Python |
|---|---|---|
| [Domain](domain.md) | `internal/core/domain/order.go` | `internal/core/domain/order.py` |
| [Inbound port](inbound_port.md) | `internal/core/application/ports/inbound/place_order.go` | `internal/core/application/ports/inbound/place_order.py` |
| [Outbound port](outbound_port.md) | `internal/core/application/ports/outbound/order_repository.go` | `internal/core/application/ports/outbound/order_repository.py` |
| [Use case](usecase.md) | `internal/core/application/usecases/place_order.go` | `internal/core/application/usecases/place_order.py` |
| [Inbound adapter](inbound_adapter.md) | `internal/adapters/inbound/http/place_order.go` | `internal/adapters/inbound/http/place_order.py` |
| [Outbound adapter](outbound_adapter.md) | `internal/adapters/outbound/postgres/order_repository.go` | `internal/adapters/outbound/postgres/order_repository.py` |
| [Bootstrap](bootstrap.md) | `bootstrap/main.go` | `bootstrap/main.py` |

## Scaling
The core and the adapters scale on their own because ports keep them apart. See [internal/ scaling](../README.md#internal).
