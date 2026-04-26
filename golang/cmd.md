# Application Entry Points
**Path:** `cmd/`

## Purpose
The `cmd` directory contains the application's entry points, also known as the composition roots. These files are responsible for initializing the application, configuring its components, and starting the execution.

## Contains
- **`cmd/[app_name]/main.go`**: The entry point for a specific application component (e.g., `cmd/api/main.go` or `cmd/worker/main.go`).

## Principles & Rules
- **Composition Root**: This is where all application layers (core, adapters, and configuration) are wired together.
- **Application Initialization**: Responsible for loading configurations and initializing dependencies (outbound adapters, use cases, and finally inbound adapters).
- **Graceful Shutdown**: Should handle signals for graceful shutdown to ensure a clean exit.
- **Minimum Logic**: Only initialization and setup logic belong here; business rules and application orchestration must be in the `internal/core` layer.
