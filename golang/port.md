# Ports
**Path:** `internal/core/application/ports/`

## Purpose
Ports are interfaces that define the contracts for interaction with the application's core system. They represent the entry and exit points for data and execution, ensuring the core remains decoupled from external infrastructure.

## Contains
- **[inbound/](inbound_port.md)**: Interfaces defining what the system can do and how other parts can use it.
- **[outbound/](outbound_port.md)**: Interfaces defining what the core needs from external systems (e.g., databases, external APIs).

## Principles & Rules
- **Interface-Based**: Ports are always defined as Go interfaces.
- **Abstraction**: Ports must be technology-agnostic (e.g., `UserRepositoryPort`, not `UserSQLRepositoryPort`).
- **Core Ownership**: Ports are defined within the core and owned by the core, even though they may be implemented outside the core (outbound) or called from outside (inbound).
