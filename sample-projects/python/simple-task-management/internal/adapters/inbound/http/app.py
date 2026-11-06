from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from internal.adapters.inbound.http.error_handlers import register_error_handlers
from internal.adapters.inbound.http.routes.board_routes import create_board_router
from internal.adapters.inbound.http.routes.task_routes import create_task_router
from internal.core.application.ports.inbound.board import (
    AddColumnUseCase,
    CreateBoardUseCase,
    GetBoardUseCase,
    ListBoardsUseCase,
)
from internal.core.application.ports.inbound.task import (
    CreateTaskUseCase,
    DeleteTaskUseCase,
    ListTasksUseCase,
    MoveTaskUseCase,
)


def create_app(
    create_board_uc: CreateBoardUseCase,
    get_board_uc: GetBoardUseCase,
    list_boards_uc: ListBoardsUseCase,
    add_column_uc: AddColumnUseCase,
    create_task_uc: CreateTaskUseCase,
    move_task_uc: MoveTaskUseCase,
    list_tasks_uc: ListTasksUseCase,
    delete_task_uc: DeleteTaskUseCase,
    title: str = "Simple Task Management API",
    version: str = "1.0.0",
) -> FastAPI:
    app = FastAPI(
        title=title,
        version=version,
        description="Kanban-style Simple Task Management System REST API",
    )

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    register_error_handlers(app)

    board_router = create_board_router(
        create_board_uc=create_board_uc,
        get_board_uc=get_board_uc,
        list_boards_uc=list_boards_uc,
        add_column_uc=add_column_uc,
        create_task_uc=create_task_uc,
        list_tasks_uc=list_tasks_uc,
    )
    task_router = create_task_router(
        move_task_uc=move_task_uc,
        delete_task_uc=delete_task_uc,
    )

    app.include_router(board_router)
    app.include_router(task_router)

    @app.get("/health", tags=["system"])
    def health():
        return {"status": "ok"}

    return app

