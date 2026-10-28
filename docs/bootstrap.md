# Application Bootstrap
**Path:** `bootstrap/`

## Purpose
The `bootstrap` directory contains the application's entry points, also known as the composition roots. These files are responsible for initializing the application, configuring its components, and starting the execution.

> **Naming:** Go projects usually call this directory `cmd/`. We use `bootstrap/` because the name works in both Go and Python and says what the code does: it builds the application and starts it. In Python, `cmd` is a standard library module, and a top-level package named `cmd` would shadow it.

## Contains
- **`bootstrap/[app_name]/[entrypoint].*`**: The entry point for a specific application component (e.g., `bootstrap/myapp/http.*`, `bootstrap/myapp/worker.*`, `bootstrap/myapp/cli.*`, ...).

## Principles & Rules
- **Composition Root**: This is where all application layers (core, adapters, configuration and etc.) are wired together.
- **Application Initialization**: Responsible for loading configurations and initializing dependencies (outbound adapters, use cases, inbound adapters and etc.).
- **Graceful Shutdown**: Should handle signals for graceful shutdown to ensure a clean exit.
- **Minimum Logic**: Only initialization and setup logic belong here; business rules and application orchestration must be in the `internal/core` layer.
