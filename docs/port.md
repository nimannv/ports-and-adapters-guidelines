# Ports
**Path:** `internal/core/application/ports/`

## Purpose
Ports are interfaces that define the contracts for interaction with the application's core system. They represent the entry and exit points for data and execution, ensuring the core remains decoupled from external infrastructure.

## Contains
- **[inbound/](inbound_port.md)**: Interfaces defining what the system can do and how other parts can use it.
- **[outbound/](outbound_port.md)**: Interfaces defining what the core needs from external systems (e.g., databases, external APIs).

## Principles & Rules
- **Interface-Based**: Ports are always defined as interfaces, together with their DTOs and errors.
- **Abstraction**: Ports must be technology-agnostic (e.g., `OrderRepository`, not `PostgresOrderRepository`).
- **Core Ownership**: Ports are defined within the core and owned by the core, even though they may be implemented outside the core (outbound) or called from outside (inbound).
- **Own Their DTOs and Errors**: Each port declares the data it accepts and returns, and the errors its callers can handle. Domain types never appear in a port signature.

## Go and Python
| Concept | Go | Python |
|---|---|---|
| Port | `interface` | `typing.Protocol` (or `abc.ABC` if you want explicit inheritance) |
| DTO | plain `struct` | `@dataclass(frozen=True)` |
| Error | sentinel `var ErrX = errors.New(...)` or error type | `class XError(Exception)` |
| Implementation check | `var _ inbound.PlaceOrder = (*placeOrder)(nil)` | type checker (`mypy` / `pyright`) |
| Context / cancellation | first parameter `ctx context.Context` | not needed (use `async def` if the app is async) |

> **Async Python:** if the application is async, make every port method `async def` and implement them with async drivers. Don't mix sync and async ports in one application.

## Scaling
Inbound ports follow the use-case split. Outbound ports are split by the capability the core needs (interface segregation), never by technology. See the [Per-Layer Guide](../README.md#per-layer-guide).
