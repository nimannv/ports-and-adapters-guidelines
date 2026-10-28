# Domain Layer
**Path:** `internal/core/domain/`

## Purpose
The `domain` layer defines and encapsulates the core business model, rules, and data structures of the application. It represents the heart of the business domain and is independent of application-specific logic or external technologies.

## Contains
- **`entities/`**: Core business objects with identity and lifecycle (e.g., `User`, `Product`).
- **`value_objects/`**: Objects representing descriptive aspects of the domain without unique identity (e.g., `EmailAddress`, `Money`).
- **`exceptions/`**: Custom error types specific to the domain failures or invalid states.
- **`services/`**: "Domain Services" encapsulating business logic that doesn't fit into a single entity or value object.

## Principles & Rules
- **Pure Business Logic**: Only contains logic essential to the business domain, free from infrastructure or application orchestration.
- **No External Dependencies**: Must not depend on anything outside the `domain` package (except standard Go libraries).
- **No Knowledge of Ports**: Completely unaware of `inbound` or `outbound` ports.
- **High Cohesion**: Entities, value objects, exceptions, and services within this layer must be tightly related to the domain concepts they represent.
