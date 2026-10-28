# Outbound Ports
**Path:** `internal/core/application/ports/outbound/`

## Purpose
Outbound ports define the external systems or capabilities that the application core needs from the outside world. This includes databases, message queues, external APIs, or file systems.

## Contains
- **Interfaces**: Definitions of methods that the core calls to interact with external systems.
- **DTOs**: Data Transfer Objects used to pass data back and forth between the core and external systems.
- **Errors**: The errors the core can react to (e.g., `ErrOrderNotFound`, `ErrDuplicateOrder`).

## Principles & Rules
- **Dependency Inversion**: The core defines its own needs as ports.
- **Technology Agnostic**: Ports define *what* is needed, not *how* it's done (e.g., `OrderRepository` instead of `PostgresOrderRepository`).
- **Named by Capability**: `OrderRepository`, `PaymentGateway`, `Notifier`, ... never by technology.
- **Small Interfaces**: Split an interface when its consumers each use only part of it (interface segregation).
- **Implementation**: These interfaces are implemented by [Outbound Adapters](outbound_adapter.md).
- **Error Handling**: Ports define the error types they can return. [Outbound adapters](outbound_adapter.md) translate technology errors into them.

## Example

```go
// internal/core/application/ports/outbound/order_repository.go
package outbound

import (
	"context"
	"errors"
)

var ErrDuplicateOrder = errors.New("order already exists")

type OrderItemRecord struct {
	ProductID string
	Quantity  int
}

type OrderRecord struct {
	ID         string
	CustomerID string
	Status     string
	Items      []OrderItemRecord
}

type OrderRepository interface {
	NextID(ctx context.Context) (string, error)
	Save(ctx context.Context, order OrderRecord) error
}
```

```python
# internal/core/application/ports/outbound/order_repository.py
from dataclasses import dataclass
from typing import Protocol


class DuplicateOrderError(Exception):
    pass


@dataclass(frozen=True)
class OrderItemRecord:
    product_id: str
    quantity: int


@dataclass(frozen=True)
class OrderRecord:
    id: str
    customer_id: str
    status: str
    items: tuple[OrderItemRecord, ...]


class OrderRepository(Protocol):
    def next_id(self) -> str: ...
    def save(self, order: OrderRecord) -> None: ...
```

## Scaling
Split by the capability the core needs, **never by technology**. See the [Per-Layer Guide](../README.md#per-layer-guide).
