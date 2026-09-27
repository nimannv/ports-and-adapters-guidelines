# Internal Packages
**Path:** `internal/`

## Purpose
The `internal` package contains the core business logic and all implementation details of the application.

## Contains
- **[core/](internal_core.md)**: The heart of the application, containing domain models and use cases.
- **[adapters/](adapters.md)**: Interactions with the external world 
    - **outbound**: databases, APIs, message brokers, etc.
    - **inbound**: HTTP, gRPC, etc.

## Principles & Rules
- **Encapsulation**: All application-specific logic must reside here.
- **Separation of Concerns**: Business logic (core) is kept strictly separate from technical details (adapters).
