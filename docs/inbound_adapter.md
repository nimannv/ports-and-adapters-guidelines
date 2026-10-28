# Inbound Adapters
**Path:** `internal/adapters/inbound/`

## Purpose
Inbound adapters are external components that drive the application's core logic. They translate external requests (e.g., HTTP, gRPC, CLI commands, AMQP messages) into inbound port DTOs and call the core via its inbound ports (implemented by usecases).

## Contains
One directory per technology:
- **http/**: HTTP handlers and route definitions.
- **grpc/**: gRPC server implementations and service handlers.
- **cli/**: Command-line interface definitions and logic.
- **amqp/**: Consumers for handling incoming messages from queues.
- ...

## Principles & Rules
- **Depends on Inbound Ports**: Inbound adapters must only depend on the interfaces, DTOs, and errors defined in `internal/core/application/ports/inbound`.
- **Input Parsing & Validation**: Responsible for input parsing (e.g., JSON decoding) and basic structural validation (e.g., required fields, formats). Business rules stay in the core.
- **Output Serialization**: Responsible for converting core responses into external formats (e.g., JSON responses, gRPC messages).
- **Error Translation**: Map inbound port errors to protocol-specific results (e.g., `InvalidOrder` → HTTP `422`, gRPC `InvalidArgument`). Unknown errors become a generic internal error.
- **Independent DTOs**: Inbound adapters should define their own DTOs (request/response models) to avoid coupling the core to the transport.
- **Technology Specific**: They are the bridge between the core and the outside world, so technology-specific details (e.g., Go `net/http`, Python FastAPI) belong here.

## Scaling

### When to Split
- **A protocol has many handlers**: one directory per technology is no longer readable.
- **Audiences differ**: public, admin, and internal callers need different auth, middleware, rate limits, or SLAs.
- **The contract is versioned**: old and new API versions must be served side by side.

### How to Split
Split in this order: **technology → scope → domain (→ version)**.

1. **Technology**: one directory per protocol (`http/`, `grpc/`, `cli/`, `amqp/`). Always the first level.
2. **Scope**: audience inside a protocol (`http/public/`, `http/admin/`, `http/internal/`). Each scope owns its own middleware and auth.
3. **Domain**: business concept inside a scope (`http/public/order/`, `http/public/billing/`), matching [usecases/](usecase.md#scaling).
4. **Version**: only when versions must coexist (`grpc/v1/`, `http/public/v2/`).

```
# 1. Flat per technology    # 2. Technology → scope → domain (→ version)
inbound/                    inbound/
├── http/            ──►    ├── http/
│   ├── order.*             │   ├── public/
│   └── billing.*           │   │   ├── order/
└── amqp/                   │   │   └── billing/
    └── billing.*           │   └── admin/
                            ├── grpc/
                            │   └── v1/
                            └── amqp/
                                └── billing/
```

Skip levels you don't need. A service with one audience and one version goes straight from technology to domain (`http/order/`).

See [Scale](../README.md#scale) for the general principles.
