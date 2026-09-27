# Core Layer
**Path:** `internal/core/`

## Purpose
The `core` layer encapsulates the heart of the application's business logic, completely decoupled from external technologies (such as web frameworks, database drivers, or external APIs). It defines *what* the application does from a business perspective.

## Contains
- **[domain/](domain.md)**: Pure business entities, value objects, domain exceptions, and domain services.
- **[application/](internal_core_application.md)**: Application-specific use cases and orchestration logic.

## Principles & Rules
- **Technology Agnostic**: This layer must not have any direct dependencies on external frameworks or infrastructure (e.g., `webservers`, `databases`, `message brokers`, ...).
- **High Cohesion, Low Coupling**: Code within the core is tightly related to the business problem, interacting with external layers strictly through defined ports.
- **Testable in Isolation**: The entire `core` layer must be unit-testable without requiring external infrastructure, using mocks for outbound ports.
