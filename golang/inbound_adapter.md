# Inbound Adapters
**Path:** `internal/adapters/inbound/`

## Purpose
Inbound adapters are external components that drive the application's core logic. They translate external requests (e.g., HTTP, gRPC, CLI commands) into internal DTOs and call the core via its inbound ports.

## Contains
- **http/**: HTTP handlers and route definitions.
- **grpc/**: gRPC server implementations and service handlers.
- **cli/**: Command-line interface definitions and logic.
- **message_queue_consumer/**: Consumers for handling incoming messages from queues.

## Principles & Rules
- **Depends on Inbound Ports**: Inbound adapters must only depend on the interfaces defined in `internal/core/application/ports/inbound`.
- **Input Parsing & Validation**: Responsible for input parsing (e.g., JSON decoding) and basic structural validation (e.g., field formats).
- **Output Serialization**: Responsible for converting core responses into external formats (e.g., JSON responses, gRPC status codes).
- **Technology Specific**: They are the bridge between the core and the outside world, so technology-specific details (e.g., `net/http` package) belong here.
