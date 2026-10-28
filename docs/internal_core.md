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
The core grows in stages (see [Growth Stages of the Core](../README.md#growth-stages-of-the-core)):

1. **Flat**: one small domain, each layer is a flat list of files.
2. **Split inside layers**: keep `domain/` and `application/` on top, group by concept inside each (`domain/order/`, `usecases/order/`).
3. **Split by module**: put the bounded context on top (`core/order/`, `core/billing/`), each with its own `domain/` and `application/`. Modules talk to each other only through ports and never import each other's `domain/`. Shared code lives in a tiny `core/shared/` kernel.
4. **Split into services**: move a module to its own repository or service.
