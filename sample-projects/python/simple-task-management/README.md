# Simple Task Management — Hexagonal Architecture (Python)

A reference implementation of the **[Simple Task Management Specification](../../sample-projects/simple-task-management.md)** using the **Ports and Adapters (Hexagonal Architecture)** pattern in Python.

The system manages Kanban boards, customizable workflow columns, and tasks, with strict business invariant enforcement including **Work-in-Progress (WIP) limits** per column. It provides dual interaction interfaces (**REST API** and **CLI**) and switchable relational persistence (**SQLite** and **PostgreSQL**).

---

## Table of Contents

- [Architectural Principles](#architectural-principles)
- [Folder Structure](#folder-structure)
- [Layers & Decoupling](#layers--decoupling)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Local Installation](#local-installation)
- [Running the System](#running-the-system)
  - [1. Running the REST API](#1-running-the-rest-api)
  - [2. Running the CLI Interface](#2-running-the-cli-interface)
  - [3. Running with Docker & Docker Compose](#3-running-with-docker--docker-compose)
- [CLI Usage Examples](#cli-usage-examples)
- [REST API Endpoints & Examples](#rest-api-endpoints--examples)
- [Switchable Persistence (SQLite / PostgreSQL)](#switchable-persistence-sqlite--postgresql)
- [Testing & Architecture Verification](#testing--architecture-verification)

---

## Architectural Principles

This project adheres strictly to the Ports and Adapters guidelines documented in this repository:

1. **Dependency Inversion**: Dependencies always point inward toward the core:
   - `core/domain`: Depends on nothing (Python standard library only).
   - `core/application/ports`: Interfaces and DTOs; depend on nothing.
   - `core/application/usecases`: Implement inbound ports; depend only on domain and ports.
   - `adapters`: Depend only on ports. Inbound adapters depend on inbound ports; outbound adapters implement outbound ports.
   - `bootstrap`: Composition root; wires components together and loads configurations.
2. **Domain Isolation**: Domain entities (`Board`, `Column`, `Task`, `Priority`) and domain exceptions contain pure business logic and invariants. Domain types never leak outside the core.
3. **Dedicated Port DTOs**: Input and output Data Transfer Objects are defined explicitly at the ports boundary to decouple transport and database models from core entities.
4. **Parity Between Interfaces**: Both the FastAPI REST API and Typer CLI interact with identical use case inbound ports.
5. **Architectural Verification**: Boundary contracts are enforced automatically via `import-linter`.

```mermaid
flowchart TD
    subgraph Driving Adapters [Inbound Adapters]
        HTTP["FastAPI REST API\n(internal/adapters/inbound/http)"]
        CLI["Typer CLI\n(internal/adapters/inbound/cli)"]
    end

    subgraph Core [Application Core]
        direction TB
        subgraph PortsIn [Inbound Ports]
            IP["CreateBoard, GetBoard, AddColumn,\nCreateTask, MoveTask, ListTasks, ..."]
        end

        subgraph Application [Use Cases]
            UC["CreateBoardService, MoveTaskService,\nCreateTaskService, ..."]
        end

        subgraph Domain [Pure Domain]
            DOM["Board, Column, Task, Priority,\nDomain Invariants & Errors"]
        end

        subgraph PortsOut [Outbound Ports]
            OP["BoardRepository, TaskRepository\nRecords & Interfaces"]
        end
    end

    subgraph Driven Adapters [Outbound Adapters]
        SQL["SQLAlchemy Adapter\n(SQLite / PostgreSQL)\n(internal/adapters/outbound/sql)"]
    end

    HTTP --> IP
    CLI --> IP
    IP -.-> UC
    UC --> DOM
    UC --> OP
    SQL -.->|implements| OP
```

---

## Folder Structure

The project layout follows the structure described in [README.md](../../README.md) and [docs/](../../docs/):

```
sample-projects/python/
├── bootstrap/                                # Application composition roots
│   ├── api/
│   │   └── main.py                           # REST API composition root & server entrypoint
│   └── cli/
│       └── main.py                           # CLI composition root & Typer entrypoint
├── configs/                                  # Configuration management
│   ├── base.yaml                             # Default configuration values
│   └── config.py                             # Configuration loader (files + env vars)
├── internal/
│   ├── core/                                 # Pure business logic and orchestration
│   │   ├── domain/                           # Core business entities and rules (Standard Library only)
│   │   │   ├── board/                        # Board & Column entities, domain errors
│   │   │   │   ├── entities.py
│   │   │   │   └── errors.py
│   │   │   └── task/                         # Task entity, Priority enum, domain errors
│   │   │       ├── entities.py
│   │   │       └── errors.py
│   │   └── application/                      # Use cases and port contracts
│   │       ├── ports/
│   │       │   ├── inbound/                  # Inbound port interfaces, DTOs & errors
│   │       │   │   ├── board/                # create_board, get_board, list_boards, add_column
│   │       │   │   └── task/                 # create_task, move_task, list_tasks, delete_task
│   │       │   └── outbound/                 # Outbound repository contracts & DTOs
│   │       │       ├── board/                # BoardRepository interface, BoardRecord
│   │       │       └── task/                 # TaskRepository interface, TaskRecord
│   │       └── usecases/                     # Inbound port implementations
│   │           ├── board/                    # CreateBoardService, AddColumnService, ...
│   │           └── task/                     # MoveTaskService (WIP checks), CreateTaskService, ...
│   └── adapters/                             # Technical infrastructure & interaction adapters
│       ├── inbound/                          # Driving adapters
│       │   ├── http/                         # FastAPI web adapter (routes, schemas, error handlers)
│       │   │   ├── app.py
│       │   │   ├── dtos.py
│       │   │   ├── error_handlers.py
│       │   │   └── routes/
│       │   └── cli/                          # Typer CLI adapter (commands & formatted output)
│       │       ├── app.py
│       │       ├── formatters.py
│       │       └── commands/
│       └── outbound/                         # Driven adapters
│           └── sql/                          # SQLAlchemy persistence adapter
│               ├── database.py               # Engine & session management (SQLite / Postgres)
│               ├── models.py                 # SQLAlchemy ORM models (BoardModel, ColumnModel, TaskModel)
│               ├── mappers.py                # Maps between ORM models and Outbound Records
│               ├── board_repository.py       # BoardRepository implementation
│               └── task_repository.py        # TaskRepository implementation
├── scripts/                                  # Utility scripts
│   ├── run_api.sh                            # Script to run FastAPI server
│   ├── run_cli.sh                            # Script to run Typer CLI
│   ├── test.sh                               # Test suite runner (pytest)
│   └── lint.sh                               # Import-linter architecture checker
├── tests/
│   ├── unit/                                 # Unit tests with in-memory fakes (zero DB/network)
│   │   ├── test_domain.py
│   │   └── test_usecases.py
│   └── integration/                          # Integration tests
│       ├── test_sql_repositories.py
│       ├── test_http_adapter.py
│       └── test_cli_adapter.py
├── .importlinter                             # Architectural boundary verification contracts
├── Dockerfile                                # Multi-purpose production Dockerfile
├── docker-compose.yml                        # Docker Compose setup with PostgreSQL & API
├── pyproject.toml                            # Project metadata and configuration
├── requirements.txt                          # Pip dependencies
├── task-cli                                  # Direct executable CLI wrapper
└── README.md                                 # This file
```

---

## Layers & Decoupling

| Layer | Path | Permitted Dependencies | Responsibility |
|---|---|---|---|
| **Domain** | `internal/core/domain` | Standard library only | Invariants, entities (`Board`, `Column`, `Task`), enums (`Priority`), domain errors. |
| **Inbound Ports** | `internal/core/application/ports/inbound` | Standard library only | Contracts for system operations (`exec` method), input/output DTOs, caller errors. |
| **Outbound Ports** | `internal/core/application/ports/outbound` | Standard library only | Storage contracts (`BoardRepository`, `TaskRepository`), persistence DTOs (`BoardRecord`, `TaskRecord`). |
| **Use Cases** | `internal/core/application/usecases` | `domain`, `ports/inbound`, `ports/outbound` | Orchestrating domain logic, enforcing WIP limits, coordinating persistence, error translation. |
| **Inbound Adapters** | `internal/adapters/inbound` | `ports/inbound` only | Translating HTTP/CLI requests into port DTOs and serializing responses. |
| **Outbound Adapters**| `internal/adapters/outbound` | `ports/outbound` only | Mapping outbound DTOs to SQLAlchemy ORM models and managing transactions. |
| **Bootstrap** | `bootstrap/` | Everything | Composition root: wires configurations, databases, use cases, and adapters. |

---

## Getting Started

### Prerequisites

- Python 3.9+
- Docker & Docker Compose (optional, for containerized run)

### Local Installation

1. Navigate to the project directory:
   ```bash
   cd sample-projects/python
   ```

2. Create a virtual environment and activate it:
   ```bash
   python3 -m venv .venv
   source .venv/bin/activate
   ```

3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

---

## Running the System

### 1. Running the REST API

To start the FastAPI REST API using SQLite (default):
```bash
./scripts/run_api.sh
# Or directly:
python -m bootstrap.api.main
```
The API is available at `http://localhost:8000`. Interactive OpenAPI documentation is accessible at `http://localhost:8000/docs`.

### 2. Running the CLI Interface

The CLI uses SQLite storage by default (`tasks.db` in current directory):
```bash
./task-cli --help
# Or using the script:
./scripts/run_cli.sh --help
```

### 3. Running with Docker & Docker Compose

To spin up the entire system with **PostgreSQL** and the **FastAPI REST API**:
```bash
docker compose up --build
```
This runs:
- `task-management-postgres`: PostgreSQL 16 database with health check and persistent volume.
- `task-management-api`: FastAPI REST API on port `8000`.

To run CLI commands against the Dockerized PostgreSQL instance:
```bash
DATABASE_URL="postgresql+psycopg2://postgres:postgrespassword@localhost:5432/task_management" ./task-cli board list
```

---

## CLI Usage Examples

The CLI provides complete parity with the specification:

```bash
# 1. Create a board (comes with default Backlog, In Progress [limit: 3], Done columns)
$ ./task-cli board create --name "Sprint 1" --description "Release MVP"
Board created:
ID:          3fa85f64-5717-4562-b3fc-2c963f66afa6
Name:        Sprint 1
Columns:     Backlog (unlimited), In Progress (limit: 3), Done (unlimited)

# 2. Add a column with a WIP limit
$ ./task-cli column add --board-id 3fa85f64-5717-4562-b3fc-2c963f66afa6 \
                        --name "QA Review" \
                        --wip-limit 2
Column 'QA Review' added with WIP limit 2.

# 3. Create a task in Backlog
$ ./task-cli task create --board-id 3fa85f64-5717-4562-b3fc-2c963f66afa6 \
                         --column "Backlog" \
                         --title "Implement Login" \
                         --priority HIGH \
                         --due-date "2026-12-01"
Task created: ID 7c9e6679-7425-40de-944b-e07fc1f90ae7 [HIGH]

# 4. Move a task to another column
$ ./task-cli task move --task-id 7c9e6679-7425-40de-944b-e07fc1f90ae7 --to-column "In Progress"
Task 7c9e6679-7425-40de-944b-e07fc1f90ae7 moved to 'In Progress'.

# 5. List tasks on a board
$ ./task-cli task list --board-id 3fa85f64-5717-4562-b3fc-2c963f66afa6
+--------------------------------------+-----------------+-------------+----------+
| ID                                   | TITLE           | COLUMN      | PRIORITY |
+--------------------------------------+-----------------+-------------+----------+
| 7c9e6679-7425-40de-944b-e07fc1f90ae7 | Implement Login | In Progress | HIGH     |
+--------------------------------------+-----------------+-------------+----------+

# 6. Delete a task
$ ./task-cli task delete --task-id 7c9e6679-7425-40de-944b-e07fc1f90ae7
Task 7c9e6679-7425-40de-944b-e07fc1f90ae7 deleted successfully.
```

---

## REST API Endpoints & Examples

### Endpoints Overview

| Method | Endpoint | Description | Status Code |
|---|---|---|---|
| `POST` | `/api/v1/boards` | Create a board | `201 Created` |
| `GET` | `/api/v1/boards` | List all boards | `200 OK` |
| `GET` | `/api/v1/boards/{boardId}` | Get board details with columns | `200 OK` / `404 Not Found` |
| `POST` | `/api/v1/boards/{boardId}/columns` | Add a column to a board | `201 Created` / `400 Bad Request` |
| `POST` | `/api/v1/boards/{boardId}/tasks` | Create a task | `201 Created` / `422 Unprocessable` |
| `GET` | `/api/v1/boards/{boardId}/tasks` | List tasks on board (`?column_id=`) | `200 OK` |
| `PATCH`| `/api/v1/tasks/{taskId}/move` | Move task to a column | `200 OK` / `409 Conflict` |
| `DELETE`| `/api/v1/tasks/{taskId}` | Delete a task | `204 No Content` / `404 Not Found` |

### cURL Examples

#### 1. Create a Board
```bash
curl -X POST http://localhost:8000/api/v1/boards \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sprint 1",
    "description": "Release MVP"
  }'
```

#### 2. Add a Column with a WIP Limit
```bash
curl -X POST http://localhost:8000/api/v1/boards/<BOARD_ID>/columns \
  -H "Content-Type: application/json" \
  -d '{
    "name": "QA Review",
    "wip_limit": 2
  }'
```

#### 3. Create a Task
```bash
curl -X POST http://localhost:8000/api/v1/boards/<BOARD_ID>/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "column": "Backlog",
    "title": "Implement Login",
    "priority": "HIGH"
  }'
```

#### 4. Move a Task to Another Column
```bash
curl -X PATCH http://localhost:8000/api/v1/tasks/<TASK_ID>/move \
  -H "Content-Type: application/json" \
  -d '{
    "target_column_id": "<TARGET_COLUMN_ID>"
  }'
```

**Success Response (`200 OK`):**
```json
{
  "task_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "board_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "column_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
  "column_name": "In Progress",
  "title": "Implement Login",
  "priority": "HIGH",
  "updated_at": "2026-11-03T13:20:00Z"
}
```

**WIP Limit Rejection Response (`409 Conflict`):**
```json
{
  "error": {
    "code": "WIP_LIMIT_EXCEEDED",
    "message": "Cannot move task to 'In Progress': column limit of 3 tasks has been reached.",
    "column_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "current_count": 3,
    "limit": 3
  }
}
```

---

## Switchable Persistence (SQLite / PostgreSQL)

Persistence engine selection is handled during application startup in `bootstrap/` via environment variables:

- **SQLite** (Default):
  ```bash
  export DB_DRIVER=sqlite
  export DATABASE_URL=sqlite:///tasks.db
  ```
- **PostgreSQL**:
  ```bash
  export DB_DRIVER=postgres
  export DATABASE_URL=postgresql+psycopg2://postgres:postgrespassword@localhost:5432/task_management
  ```

Because outbound ports (`BoardRepository`, `TaskRepository`) define contracts by capability rather than technology, switching the underlying database engine leaves all application use cases and domain logic 100% untouched.

---

## Testing & Architecture Verification

### Running the Test Suite
The test suite covers pure unit tests (domain entities and use cases using in-memory fakes) and integration tests (SQLAlchemy with SQLite, FastAPI TestClient, and Typer CliRunner):
```bash
./scripts/test.sh
```

### Verifying Hexagonal Architecture Boundaries
Architectural rules are verified with `import-linter` according to `.importlinter`:
```bash
./scripts/lint.sh
```

**Contracts Enforced:**
1. `core.domain` depends on nothing (standard library only).
2. `core.application.ports.inbound` depends on nothing.
3. `core.application.ports.outbound` depends on nothing.
4. `core` never imports `adapters` or `bootstrap`.
5. `adapters` never import `domain` or `usecases`.
6. `inbound adapters` never import `outbound adapters` or `outbound ports`.
7. `outbound adapters` never import `inbound adapters` or `inbound ports`.

