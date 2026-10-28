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
- **Error Translation**: Map inbound port errors to protocol-specific results (e.g., `ErrInvalidOrder` → HTTP `422`, gRPC `InvalidArgument`). Unknown errors become a generic internal error.
- **Independent DTOs**: Inbound adapters should define their own DTOs (request/response models) to avoid coupling the core to the transport.
- **Technology Specific**: They are the bridge between the core and the outside world, so technology-specific details (e.g., Go `net/http`, Python FastAPI) belong here.

## Example

```go
// internal/adapters/inbound/http/place_order.go
package http

import (
	"encoding/json"
	"errors"
	stdhttp "net/http" // aliased: this package is also named http

	"example.com/shop/internal/core/application/ports/inbound"
)

type placeOrderRequest struct {
	CustomerID string `json:"customer_id"`
	Items      []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

type placeOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type PlaceOrderHandler struct {
	placeOrder inbound.PlaceOrder
}

func NewPlaceOrderHandler(placeOrder inbound.PlaceOrder) *PlaceOrderHandler {
	return &PlaceOrderHandler{placeOrder: placeOrder}
}

func (h *PlaceOrderHandler) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req placeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		stdhttp.Error(w, "malformed JSON", stdhttp.StatusBadRequest)
		return
	}

	in := inbound.PlaceOrderInput{CustomerID: req.CustomerID}
	for _, it := range req.Items {
		in.Items = append(in.Items, inbound.PlaceOrderItem{ProductID: it.ProductID, Quantity: it.Quantity})
	}

	out, err := h.placeOrder.Exec(r.Context(), in)
	switch {
	case errors.Is(err, inbound.ErrInvalidOrder):
		stdhttp.Error(w, err.Error(), stdhttp.StatusUnprocessableEntity)
		return
	case err != nil:
		stdhttp.Error(w, "internal error", stdhttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stdhttp.StatusCreated)
	_ = json.NewEncoder(w).Encode(placeOrderResponse{OrderID: out.OrderID, Status: out.Status})
}
```

```python
# internal/adapters/inbound/http/place_order.py
from fastapi import APIRouter, HTTPException, status
from pydantic import BaseModel

from internal.core.application.ports.inbound.place_order import (
    InvalidOrderError,
    PlaceOrder,
    PlaceOrderInput,
    PlaceOrderItem,
)


class PlaceOrderItemRequest(BaseModel):
    product_id: str
    quantity: int


class PlaceOrderRequest(BaseModel):
    customer_id: str
    items: list[PlaceOrderItemRequest]


class PlaceOrderResponse(BaseModel):
    order_id: str
    status: str


def new_place_order_router(place_order: PlaceOrder) -> APIRouter:
    router = APIRouter()

    @router.post("/orders", status_code=status.HTTP_201_CREATED)
    def post_order(req: PlaceOrderRequest) -> PlaceOrderResponse:
        try:
            out = place_order.exec(
                PlaceOrderInput(
                    customer_id=req.customer_id,
                    items=tuple(PlaceOrderItem(i.product_id, i.quantity) for i in req.items),
                )
            )
        except InvalidOrderError as e:
            raise HTTPException(status.HTTP_422_UNPROCESSABLE_ENTITY, str(e)) from e
        return PlaceOrderResponse(order_id=out.order_id, status=out.status)

    return router
```

## Scaling
Split by technology → scope → domain (→ version) when a protocol has many handlers, or audiences need different auth or middleware: `http/public/order/`, `http/admin/`, `grpc/v1/`. See the [Per-Layer Guide](../README.md#per-layer-guide).
