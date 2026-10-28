# Domain Layer
**Path:** `internal/core/domain/`

## Purpose
The `domain` layer defines and encapsulates the core business model, rules, and data structures of the application. It represents the heart of the business domain and is independent of application-specific logic or external technologies.

The modeling style is not fixed. Pick the one that fits the complexity of the business logic: anemic, rich domain model, full DDD, state machine, ... (see [Domain Development](../README.md#domain-development)).

## Contains
- **Entities**: Core business objects with identity and lifecycle (e.g., `Order`, `Customer`).
- **Value Objects**: Immutable objects describing the domain without unique identity (e.g., `Money`, `EmailAddress`).
- **Domain Errors**: Errors for domain failures or invalid states (Go: error values or types; Python: exception classes).
- **Domain Services**: Business logic that doesn't fit into a single entity or value object.
- **Aggregates, Domain Events, ...**: When a richer modeling style (e.g., DDD) is used.

> **Group by business concept, not by kind.** Keep `Order`, its value objects, errors, and rules together in `order.*` (or `order/` once it grows). Don't create technical buckets such as `entities/`, `value_objects/`, `exceptions/`, or `services/`.

## Principles & Rules
- **Pure Business Logic**: Only contains logic essential to the business domain, free from infrastructure or application orchestration.
- **No External Dependencies**: Must not depend on anything outside the `domain` package, except the language standard library. No ORMs, frameworks, or serialization tags.
- **No Knowledge of Ports**: Completely unaware of use cases and of `inbound` or `outbound` ports.
- **Protects Its Invariants**: Constructors/factories validate input, so an invalid entity can't exist.
- **High Cohesion**: Entities, value objects, errors, and services within this layer must be tightly related to the domain concepts they represent.

## Scaling

### When to Split
- Entities form separate groups with their own rules and invariants.
- Groups change for different reasons or at different rates.
- The flat directory is hard to read at a glance (roughly 10–15 files).

### How to Split
Split by **aggregate / business concept**. Each group holds its entities, value objects, errors, and services together.

```
# 1. Flat                 # 2. Split by concept
domain/                   domain/
├── order.*        ──►    ├── order/
├── order_item.*          │   ├── order.*
├── customer.*            │   ├── order_item.*
└── address.*             │   └── errors.*
                          └── customer/
                              ├── customer.*
                              └── address.*
```

- Concepts that reference each other a lot belong in the same group. If two groups import each other, merge them or extract the shared concept.
- When the same word means different things in different parts of the business (e.g. `Customer` in sales vs. billing), you have separate bounded contexts. Split the whole core instead (see [Core Scaling, Stage 3](internal_core.md#stage-3-split-by-module)).

See [Scale](../README.md#scale) for the general principles.
