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
| Port | `interface` | `typing.Protocol` (or `abc.ABC` for explicit inheritance) |
| DTO | plain `struct` | frozen `dataclass` |
| Error | error value or error type | exception class |
| Implementation check | compile-time interface assertion | static type checker (`mypy` / `pyright`) |
| Cancellation / deadlines | `context.Context` as first parameter | not needed (use `async` if the app is async) |

> **Async Python:** if the application is async, make every port method `async` and implement them with async drivers. Don't mix sync and async ports in one application.

## Scaling
- **[inbound/](inbound_port.md#scaling)**: follows the use-case split.
- **[outbound/](outbound_port.md#scaling)**: by capability the core needs, never by technology.
