# Use Cases
**Path:** `internal/core/application/usecases/`

## Purpose
The `usecases` package contains concrete implementations of the [inbound ports](./inbound_port.md). It defines *how* the application performs its specific tasks (application layer logic) by orchestrating domain logic, using outbound ports and other usecases.

## Contains
- **Usecase Implementations**: One Go or Python file per use case (e.g., `place_order.*`), implementing an interface defined in [inbound ports](./inbound_port.md).

## Principles & Rules

- **Dependency Direction**:
    - **Depends on Domain**: Directly uses entities, value objects, and domain services.
    - **Depends on Outbound Ports**: Depends on interfaces defined in `application/ports/outbound` to interact with external systems.
    - **Depends on other Use Cases**: Depends on interfaces defined in `application/ports/inbound` to use other usecases (never on their concrete implementations).

- **Orchestration Layer**:
    - Coordinates domain logic
    - Interactions with external systems via outbound ports
    - Interactions with other usecases

- **DTO and Error Mapping**:
    - Maps inbound port DTOs to domain objects, and domain objects to outbound port DTOs (and back).
    - Translates domain errors into the errors declared by the inbound port.

- **Transactional Consistency**:
    - Often serves as the boundary for database transactions. Express the transaction through an outbound port (e.g., a `UnitOfWork` / `TxManager`), never through a database driver.

- **Dependency Injection**:
    - Receives its dependencies (outbound ports, other use cases) through its constructor. They are wired in [bootstrap](bootstrap.md).

## Example

```go
// internal/core/application/usecases/place_order.go
package usecases

import (
	"context"
	"fmt"

	"example.com/shop/internal/core/application/ports/inbound"
	"example.com/shop/internal/core/application/ports/outbound"
	"example.com/shop/internal/core/domain"
)

type placeOrder struct {
	orders outbound.OrderRepository
}

// Compile-time check that placeOrder implements the inbound port.
var _ inbound.PlaceOrder = (*placeOrder)(nil)

func NewPlaceOrder(orders outbound.OrderRepository) inbound.PlaceOrder {
	return &placeOrder{orders: orders}
}

func (uc *placeOrder) Exec(ctx context.Context, in inbound.PlaceOrderInput) (inbound.PlaceOrderOutput, error) {
	id, err := uc.orders.NextID(ctx)
	if err != nil {
		return inbound.PlaceOrderOutput{}, fmt.Errorf("next order id: %w", err)
	}

	items := make([]domain.OrderItem, 0, len(in.Items))
	for _, it := range in.Items {
		items = append(items, domain.OrderItem{ProductID: it.ProductID, Quantity: it.Quantity})
	}

	order, err := domain.NewOrder(id, in.CustomerID, items)
	if err != nil {
		return inbound.PlaceOrderOutput{}, fmt.Errorf("%w: %w", inbound.ErrInvalidOrder, err)
	}

	if err := uc.orders.Save(ctx, toOrderRecord(order)); err != nil {
		return inbound.PlaceOrderOutput{}, fmt.Errorf("save order: %w", err)
	}

	return inbound.PlaceOrderOutput{OrderID: order.ID, Status: string(order.Status)}, nil
}

func toOrderRecord(o *domain.Order) outbound.OrderRecord {
	items := make([]outbound.OrderItemRecord, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, outbound.OrderItemRecord{ProductID: it.ProductID, Quantity: it.Quantity})
	}
	return outbound.OrderRecord{ID: o.ID, CustomerID: o.CustomerID, Status: string(o.Status), Items: items}
}
```

```python
# internal/core/application/usecases/place_order.py
from internal.core.application.ports.inbound.place_order import (
    InvalidOrderError,
    PlaceOrderInput,
    PlaceOrderOutput,
)
from internal.core.application.ports.outbound.order_repository import (
    OrderItemRecord,
    OrderRecord,
    OrderRepository,
)
from internal.core.domain.order import DomainError, Order, OrderItem


class PlaceOrderUseCase:
    """Implements the PlaceOrder inbound port."""

    def __init__(self, orders: OrderRepository) -> None:
        self._orders = orders

    def exec(self, data: PlaceOrderInput) -> PlaceOrderOutput:
        try:
            order = Order.new(
                order_id=self._orders.next_id(),
                customer_id=data.customer_id,
                items=[OrderItem(i.product_id, i.quantity) for i in data.items],
            )
        except DomainError as e:
            raise InvalidOrderError(str(e)) from e

        self._orders.save(_to_order_record(order))
        return PlaceOrderOutput(order_id=order.id, status=order.status.value)


def _to_order_record(order: Order) -> OrderRecord:
    return OrderRecord(
        id=order.id,
        customer_id=order.customer_id,
        status=order.status.value,
        items=tuple(OrderItemRecord(i.product_id, i.quantity) for i in order.items),
    )
```

## Testing
Use cases are tested in isolation with in-memory fakes of the outbound ports. No database, no network.

```go
// internal/core/application/usecases/place_order_test.go
package usecases_test

import (
	"context"
	"errors"
	"testing"

	"example.com/shop/internal/core/application/ports/inbound"
	"example.com/shop/internal/core/application/ports/outbound"
	"example.com/shop/internal/core/application/usecases"
)

type fakeOrders struct{ saved []outbound.OrderRecord }

func (f *fakeOrders) NextID(context.Context) (string, error) { return "order-1", nil }

func (f *fakeOrders) Save(_ context.Context, o outbound.OrderRecord) error {
	f.saved = append(f.saved, o)
	return nil
}

func TestPlaceOrder_RejectsEmptyOrder(t *testing.T) {
	uc := usecases.NewPlaceOrder(&fakeOrders{})
	_, err := uc.Exec(context.Background(), inbound.PlaceOrderInput{CustomerID: "c-1"})
	if !errors.Is(err, inbound.ErrInvalidOrder) {
		t.Fatalf("want ErrInvalidOrder, got %v", err)
	}
}
```

```python
# internal/core/application/usecases/test_place_order.py
import pytest

from internal.core.application.ports.inbound.place_order import InvalidOrderError, PlaceOrderInput
from internal.core.application.ports.outbound.order_repository import OrderRecord
from internal.core.application.usecases.place_order import PlaceOrderUseCase


class FakeOrders:
    def __init__(self) -> None:
        self.saved: list[OrderRecord] = []

    def next_id(self) -> str:
        return "order-1"

    def save(self, order: OrderRecord) -> None:
        self.saved.append(order)


def test_rejects_empty_order() -> None:
    uc = PlaceOrderUseCase(FakeOrders())
    with pytest.raises(InvalidOrderError):
        uc.exec(PlaceOrderInput(customer_id="c-1", items=()))
```

## Scaling
Split by feature, then by actor (or by command/query if reads and writes differ a lot): `usecases/order/place_order.*`, `usecases/admin/...`. Inbound ports follow the same split. See the [Per-Layer Guide](../README.md#per-layer-guide).
