# Domain Layer
**Path:** `internal/core/domain/`

## Purpose
The `domain` layer defines and encapsulates the core business model, rules, and data structures of the application. It represents the heart of the business domain and is independent of application-specific logic or external technologies.

The modeling style is not fixed. Pick the one that fits the complexity of the business logic: anemic, rich domain model, full DDD, state machine, ... (see [Domain Development](../README.md#domain-development)).

## Contains
- **Entities**: Core business objects with identity and lifecycle (e.g., `Order`, `Customer`).
- **Value Objects**: Immutable objects describing the domain without unique identity (e.g., `Money`, `EmailAddress`, `OrderItem`).
- **Domain Errors**: Errors for domain failures or invalid states (Go: sentinel errors or error types; Python: exception classes).
- **Domain Services**: Business logic that doesn't fit into a single entity or value object.
- **Aggregates, Domain Events, ...**: When a richer modeling style (e.g., DDD) is used.

> **Group by business concept, not by kind.** Keep `Order`, its value objects, errors, and rules together in `order.*` (or `order/` once it grows). Don't create technical buckets such as `entities/`, `value_objects/`, `exceptions/`, or `services/`.

## Principles & Rules
- **Pure Business Logic**: Only contains logic essential to the business domain, free from infrastructure or application orchestration.
- **No External Dependencies**: Must not depend on anything outside the `domain` package, except the language standard library (e.g., Go `time`, `errors`; Python `dataclasses`, `enum`, `decimal`). No ORMs, frameworks, or serialization tags.
- **No Knowledge of Ports**: Completely unaware of use cases and of `inbound` or `outbound` ports.
- **Protects Its Invariants**: Constructors/factories validate input, so an invalid entity can't exist.
- **High Cohesion**: Entities, value objects, errors, and services within this layer must be tightly related to the domain concepts they represent.

## Example
A rich `Order` entity that enforces its own rules.

```go
// internal/core/domain/order.go
package domain

import "errors"

var (
	ErrEmptyOrder      = errors.New("order must contain at least one item")
	ErrInvalidQuantity = errors.New("item quantity must be positive")
)

type OrderStatus string

const OrderStatusPending OrderStatus = "pending"

// OrderItem is a value object.
type OrderItem struct {
	ProductID string
	Quantity  int
}

// Order is an entity.
type Order struct {
	ID         string
	CustomerID string
	Items      []OrderItem
	Status     OrderStatus
}

func NewOrder(id, customerID string, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	for _, it := range items {
		if it.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
	}
	return &Order{ID: id, CustomerID: customerID, Items: items, Status: OrderStatusPending}, nil
}
```

```python
# internal/core/domain/order.py
from dataclasses import dataclass
from enum import Enum


class DomainError(Exception):
    """Base class for all domain errors."""


class EmptyOrderError(DomainError):
    pass


class InvalidQuantityError(DomainError):
    pass


class OrderStatus(str, Enum):
    PENDING = "pending"


@dataclass(frozen=True)
class OrderItem:
    """Value object."""
    product_id: str
    quantity: int


@dataclass
class Order:
    """Entity."""
    id: str
    customer_id: str
    items: list[OrderItem]
    status: OrderStatus

    @classmethod
    def new(cls, order_id: str, customer_id: str, items: list[OrderItem]) -> "Order":
        if not items:
            raise EmptyOrderError("order must contain at least one item")
        if any(item.quantity <= 0 for item in items):
            raise InvalidQuantityError("item quantity must be positive")
        return cls(order_id, customer_id, list(items), OrderStatus.PENDING)
```

## Scaling
Split when entities form separate groups with their own rules and change for different reasons. Split by aggregate or business concept: `domain/order/`, `domain/customer/`. See the [Per-Layer Guide](../README.md#per-layer-guide).
