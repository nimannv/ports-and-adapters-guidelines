# Application Bootstrap
**Path:** `bootstrap/`

## Purpose
The `bootstrap` directory contains the application's entry points, also known as the composition roots. These files are responsible for initializing the application, configuring its components, and starting the execution.

> **Naming:** Go projects usually call this directory `cmd/`. We use `bootstrap/` because the name works in both Go and Python and says what the code does: it builds the application and starts it. In Python, `cmd` is a standard library module, and a top-level package named `cmd` would shadow it.

## Contains
- **`main.*`**: One entry point per program (e.g., `bootstrap/main.*`, `bootstrap/api/main.*`, `bootstrap/order/worker/main.*`). See [Scaling](#scaling) for how the layout grows.

## Principles & Rules
- **Composition Root**: This is where all application layers (core, adapters, configuration, etc.) are wired together. It is the only place that knows about concrete implementations.
- **Application Initialization**: Loads [configurations](configs.md) and initializes dependencies in order: outbound adapters → use cases → inbound adapters.
- **Graceful Shutdown**: Handles signals (`SIGINT`, `SIGTERM`) for graceful shutdown to ensure a clean exit.
- **Minimum Logic**: Only initialization and setup logic belong here; business rules and application orchestration must be in the `internal/core` layer.
- **Wire Only What You Need**: Each entry point wires only the use cases and adapters it uses.

## Scaling
Start with one entry point. Each entry point stays a thin composition root. As the codebase grows, split by **scope** (runtime), then by **domain**.

```
# 1. Single entry point     # 2. Split by scope (runtime)    # 3. Split by domain, then scope
bootstrap/                  bootstrap/                       bootstrap/
└── main.*                  ├── api/                         ├── order/
                            │   └── main.*                   │   ├── api/
                            ├── worker/                      │   │   └── main.*
                            │   └── main.*                   │   └── worker/
                            ├── cli/                         │       └── main.*
                            │   └── main.*                   ├── billing/
                            └── migrate/                     │   ├── api/
                                └── main.*                   │   └── worker/
                                                             └── admin/
                                                                 └── cli/
                                                                     └── main.*
```

- **Split by scope** when one program has to start different kinds of work (serve HTTP, consume a queue, run a job). These need to be deployed, scaled, or restarted separately.
- **Split by domain** when separate business modules get their own deployments, release cycles, or owning teams. Each one wires only the use cases and adapters it needs.

See [Scale](../README.md#scale) for the general principles.
