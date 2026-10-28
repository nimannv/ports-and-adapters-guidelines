# Outbound Adapters
**Path:** `internal/adapters/outbound/`

## Purpose
**Outbound adapters** implement the [outbound ports](./outbound_port.md) defined by the application's core. They fulfill the core's needs by interacting with external systems, such as databases, APIs, or message queues.

## Contains
One directory per technology or external system:
- **postgres/**, **mongo/**, **redis/**: Implementations of data persistence outbound ports.
- **stripe/**, **sendgrid/**, ...: Implementations of outbound ports for calling third-party services.
- **amqp/**: Publishers for outgoing messages.
- ...

## Principles & Rules
- **Implements Outbound Ports**: Outbound adapters must implement the [outbound ports](./outbound_port.md) defined by the application's core.
- **DTO Mapping**: They use the outbound port DTOs and map them to the external systems' formats (HTTP messages, ORM models, SQL rows, etc.).
- **Technology Specific**: They are the bridge to specific external systems and contain technology-specific logic (e.g., SQL queries, API clients).
- **Error Translation**: Responsible for translating technology-specific errors (e.g., unique-constraint violations, timeouts) into the error types defined by the outbound port, so the core can understand them.
- **No Business Logic**: Their primary role is translation and interaction, not containing business rules.

## Scaling

### When to Split
- **One adapter serves several domains**: e.g. a single Postgres adapter implements repositories for orders and billing.
- **Several external systems**: a new datastore or third-party service joins.

### How to Split
Split in this order: **technology / external system → domain**.

1. **Technology / external system**: one directory per datastore or provider (`postgres/`, `redis/`, `stripe/`). Always the first level.
2. **Domain**: inside a technology, one directory per business concept (`postgres/order/`, `postgres/billing/`), matching the [outbound ports](outbound_port.md#scaling) they implement.

```
# 1. Flat per technology       # 2. Technology → domain
outbound/                      outbound/
├── postgres/           ──►    ├── postgres/
│   ├── order_repo.*           │   ├── order/
│   └── invoice_repo.*         │   └── billing/
└── stripe/                    ├── redis/
    └── payment_gateway.*      │   └── order/
                               └── stripe/
```

- Third-party clients for one provider (e.g. `stripe/`) usually stay flat; they rarely need a domain level.
- Shared technology code (connection pools, migrations helpers, retry logic) lives at the technology level (`postgres/`), not in the core.

See [Scale](../README.md#scale) for the general principles.
