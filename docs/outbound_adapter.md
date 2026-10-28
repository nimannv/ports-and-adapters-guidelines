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

## Example

```go
// internal/adapters/outbound/postgres/order_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"example.com/shop/internal/core/application/ports/outbound"
)

type OrderRepository struct {
	db *sql.DB
}

var _ outbound.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) NextID(context.Context) (string, error) {
	return uuid.NewString(), nil
}

func (r *OrderRepository) Save(ctx context.Context, o outbound.OrderRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, status) VALUES ($1, $2, $3)`,
		o.ID, o.CustomerID, o.Status)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return outbound.ErrDuplicateOrder
	}
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	for _, it := range o.Items {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
			o.ID, it.ProductID, it.Quantity); err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}
	return tx.Commit()
}
```

```python
# internal/adapters/outbound/postgres/order_repository.py
import uuid

import psycopg
from psycopg import errors as pg_errors

from internal.core.application.ports.outbound.order_repository import (
    DuplicateOrderError,
    OrderRecord,
)


class PostgresOrderRepository:
    """Implements the OrderRepository outbound port."""

    def __init__(self, conn: psycopg.Connection) -> None:
        self._conn = conn

    def next_id(self) -> str:
        return str(uuid.uuid4())

    def save(self, order: OrderRecord) -> None:
        try:
            with self._conn.transaction(), self._conn.cursor() as cur:
                cur.execute(
                    "INSERT INTO orders (id, customer_id, status) VALUES (%s, %s, %s)",
                    (order.id, order.customer_id, order.status),
                )
                cur.executemany(
                    "INSERT INTO order_items (order_id, product_id, quantity) VALUES (%s, %s, %s)",
                    [(order.id, i.product_id, i.quantity) for i in order.items],
                )
        except pg_errors.UniqueViolation as e:
            raise DuplicateOrderError(order.id) from e
```

## Scaling
Split by technology / external system, then by domain, when one adapter implements ports for several domains or talks to several external systems: `postgres/order/`, `stripe/`. See the [Per-Layer Guide](../README.md#per-layer-guide).
