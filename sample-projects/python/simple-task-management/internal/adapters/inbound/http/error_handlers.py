from fastapi import Request
from fastapi.responses import JSONResponse

from internal.core.application.ports.inbound.board import (
    BoardNotFoundError,
    CreateBoardError,
    InboundDuplicateColumnNameError,
    InboundInvalidWipLimitError,
)
from internal.core.application.ports.inbound.task import (
    InboundBoardNotFoundError,
    InboundColumnNotFoundError,
    InboundDeleteTaskNotFoundError,
    InboundInvalidTaskError,
    InboundInvalidTaskMovementError,
    InboundListTasksBoardNotFoundError,
    InboundMoveColumnNotFoundError,
    InboundMoveWipLimitExceededError,
    InboundTaskNotFoundError,
    InboundWipLimitExceededError,
)


def register_error_handlers(app):
    @app.exception_handler(BoardNotFoundError)
    @app.exception_handler(InboundBoardNotFoundError)
    @app.exception_handler(InboundListTasksBoardNotFoundError)
    async def board_not_found_handler(request: Request, exc):
        return JSONResponse(
            status_code=404,
            content={"error": {"code": "BOARD_NOT_FOUND", "message": str(exc)}},
        )

    @app.exception_handler(InboundColumnNotFoundError)
    @app.exception_handler(InboundMoveColumnNotFoundError)
    async def column_not_found_handler(request: Request, exc):
        return JSONResponse(
            status_code=404,
            content={"error": {"code": "COLUMN_NOT_FOUND", "message": str(exc)}},
        )

    @app.exception_handler(InboundTaskNotFoundError)
    @app.exception_handler(InboundDeleteTaskNotFoundError)
    async def task_not_found_handler(request: Request, exc):
        return JSONResponse(
            status_code=404,
            content={"error": {"code": "TASK_NOT_FOUND", "message": str(exc)}},
        )

    @app.exception_handler(InboundDuplicateColumnNameError)
    async def duplicate_column_handler(request: Request, exc: InboundDuplicateColumnNameError):
        return JSONResponse(
            status_code=400,
            content={"error": {"code": exc.code, "message": str(exc)}},
        )

    @app.exception_handler(InboundInvalidWipLimitError)
    async def invalid_wip_limit_handler(request: Request, exc: InboundInvalidWipLimitError):
        return JSONResponse(
            status_code=400,
            content={"error": {"code": exc.code, "message": str(exc)}},
        )

    @app.exception_handler(InboundWipLimitExceededError)
    @app.exception_handler(InboundMoveWipLimitExceededError)
    async def wip_limit_exceeded_handler(request: Request, exc):
        return JSONResponse(
            status_code=409,
            content={
                "error": {
                    "code": "WIP_LIMIT_EXCEEDED",
                    "message": str(exc),
                    "column_id": exc.column_id,
                    "current_count": exc.current_count,
                    "limit": exc.limit,
                }
            },
        )

    @app.exception_handler(InboundInvalidTaskError)
    async def invalid_task_handler(request: Request, exc: InboundInvalidTaskError):
        return JSONResponse(
            status_code=422,
            content={"error": {"code": exc.code, "message": str(exc)}},
        )

    @app.exception_handler(InboundInvalidTaskMovementError)
    async def invalid_task_movement_handler(request: Request, exc: InboundInvalidTaskMovementError):
        return JSONResponse(
            status_code=400,
            content={"error": {"code": exc.code, "message": str(exc)}},
        )

    @app.exception_handler(CreateBoardError)
    async def create_board_error_handler(request: Request, exc: CreateBoardError):
        return JSONResponse(
            status_code=400,
            content={"error": {"code": exc.code, "message": str(exc)}},
        )

