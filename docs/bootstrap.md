# Application Bootstrap
**Path:** `bootstrap/`

## Purpose
The `bootstrap` directory contains the application's entry points, also known as the composition roots. These files are responsible for initializing the application, configuring its components, and starting the execution.

> **Naming:** Go projects usually call this directory `cmd/`. We use `bootstrap/` because the name works in both Go and Python and says what the code does: it builds the application and starts it. In Python, `cmd` is a standard library module, and a top-level package named `cmd` would shadow it.

## Contains
The layout grows with the codebase (see [bootstrap/ scaling](../README.md#bootstrap)):

| Stage | Layout | Run (Go) | Run (Python) |
|---|---|---|---|
| Single entry point | `bootstrap/main.*` | `go run ./bootstrap` | `python -m bootstrap.main` |
| Split by scope (runtime) | `bootstrap/api/main.*`, `bootstrap/worker/main.*`, ... | `go run ./bootstrap/api` | `python -m bootstrap.api.main` |
| Split by domain, then scope | `bootstrap/order/api/main.*`, ... | `go run ./bootstrap/order/api` | `python -m bootstrap.order.api.main` |

## Principles & Rules
- **Composition Root**: This is where all application layers (core, adapters, configuration, etc.) are wired together. It is the only place that knows about concrete implementations.
- **Application Initialization**: Responsible for loading [configurations](configs.md) and initializing dependencies in order: outbound adapters → use cases → inbound adapters.
- **Graceful Shutdown**: Should handle signals (`SIGINT`, `SIGTERM`) for graceful shutdown to ensure a clean exit.
- **Minimum Logic**: Only initialization and setup logic belong here; business rules and application orchestration must be in the `internal/core` layer.
- **Wire Only What You Need**: Each entry point wires only the use cases and adapters it uses.

## Example

```go
// bootstrap/main.go
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	httpin "example.com/shop/internal/adapters/inbound/http"
	"example.com/shop/internal/adapters/outbound/postgres"
	"example.com/shop/internal/core/application/usecases"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL")) // or load from configs/
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// outbound adapters → use cases → inbound adapters
	orders := postgres.NewOrderRepository(db)
	placeOrder := usecases.NewPlaceOrder(orders)

	mux := http.NewServeMux()
	mux.Handle("POST /orders", httpin.NewPlaceOrderHandler(placeOrder))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
```

```python
# bootstrap/main.py
import os

import psycopg
import uvicorn
from fastapi import FastAPI

from internal.adapters.inbound.http.place_order import new_place_order_router
from internal.adapters.outbound.postgres.order_repository import PostgresOrderRepository
from internal.core.application.usecases.place_order import PlaceOrderUseCase


def build_app(conn: psycopg.Connection) -> FastAPI:
    # outbound adapters → use cases → inbound adapters
    orders = PostgresOrderRepository(conn)
    place_order = PlaceOrderUseCase(orders)

    app = FastAPI()
    app.include_router(new_place_order_router(place_order))
    return app


def main() -> None:
    with psycopg.connect(os.environ["DATABASE_URL"]) as conn:  # or load from configs/
        # uvicorn handles SIGINT/SIGTERM and shuts down gracefully.
        uvicorn.run(build_app(conn), host="0.0.0.0", port=8080)


if __name__ == "__main__":
    main()
```
