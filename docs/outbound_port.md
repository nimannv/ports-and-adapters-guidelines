# Outbound Ports
**Path:** `internal/core/application/ports/outbound/`

## Purpose
Outbound ports define the external systems or capabilities that the application core needs from the outside world. This includes databases, message queues, external APIs, or file systems.

## Contains
- **Interfaces**: Definitions of methods that the core calls to interact with external systems.
- **DTOs**: Data Transfer Objects used to pass data back and forth between the core and external systems.
- **Errors**: The errors the core can react to (e.g., `OrderNotFound`, `DuplicateOrder`).

## Principles & Rules
- **Dependency Inversion**: The core defines its own needs as ports.
- **Technology Agnostic**: Ports define *what* is needed, not *how* it's done (e.g., `OrderRepository` instead of `PostgresOrderRepository`).
- **Named by Capability**: `OrderRepository`, `PaymentGateway`, `Notifier`, ... never by technology.
- **Implementation**: These interfaces are implemented by [Outbound Adapters](outbound_adapter.md).
- **Error Handling**: Ports define the error types they can return. [Outbound adapters](outbound_adapter.md) translate technology errors into them.

## Scaling

### When to Split
- **An interface grows too wide**: its consumers each use only part of it (interface segregation). For example, most use cases only read orders and a few write them.
- **The directory grows**: many ports for unrelated concepts sit side by side.

### How to Split
- **Split interfaces by capability**: `OrderReader` and `OrderWriter` instead of one large `OrderRepository`, or `PaymentGateway` and `RefundGateway`. One adapter can still implement several of them.
- **Group files by concept**: `ports/outbound/order/`, `ports/outbound/payment/`, following the same concepts as [domain/](domain.md#scaling).
- **Never by technology**: no `ports/outbound/postgres/` or `ports/outbound/redis/`. Technology belongs to [outbound adapters](outbound_adapter.md#scaling).

```
ports/outbound/                 ports/outbound/
├── order_repository.*   ──►    ├── order/
├── payment_gateway.*           │   ├── order_reader.*
└── notifier.*                  │   └── order_writer.*
                                ├── payment/
                                │   └── payment_gateway.*
                                └── notification/
                                    └── notifier.*
```

See [Scale](../README.md#scale) for the general principles.
