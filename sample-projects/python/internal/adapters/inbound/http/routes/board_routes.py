from typing import List, Optional
from fastapi import APIRouter, Depends, Query, status

from internal.adapters.inbound.http.dtos import (
    AddColumnRequest,
    BoardResponse,
    ColumnResponse,
    CreateBoardRequest,
    CreateTaskRequest,
    TaskResponse,
)
from internal.core.application.ports.inbound.board import (
    AddColumnInput,
    AddColumnUseCase,
    CreateBoardInput,
    CreateBoardUseCase,
    GetBoardInput,
    GetBoardUseCase,
    ListBoardsInput,
    ListBoardsUseCase,
)
from internal.core.application.ports.inbound.task import (
    CreateTaskInput,
    CreateTaskUseCase,
    ListTasksInput,
    ListTasksUseCase,
)


def create_board_router(
    create_board_uc: CreateBoardUseCase,
    get_board_uc: GetBoardUseCase,
    list_boards_uc: ListBoardsUseCase,
    add_column_uc: AddColumnUseCase,
    create_task_uc: CreateTaskUseCase,
    list_tasks_uc: ListTasksUseCase,
) -> APIRouter:
    router = APIRouter(prefix="/api/v1/boards", tags=["boards"])

    @router.post("", response_model=BoardResponse, status_code=status.HTTP_201_CREATED)
    def create_board(payload: CreateBoardRequest):
        custom_cols = None
        if payload.custom_columns:
            custom_cols = [c.dict() for c in payload.custom_columns]

        res = create_board_uc.exec(
            CreateBoardInput(
                name=payload.name,
                description=payload.description,
                custom_columns=custom_cols,
            )
        )
        return BoardResponse(
            id=res.id,
            name=res.name,
            description=res.description,
            columns=[
                ColumnResponse(
                    id=c.id,
                    name=c.name,
                    order_index=c.order_index,
                    wip_limit=c.wip_limit,
                )
                for c in res.columns
            ],
            created_at=res.created_at,
            updated_at=res.updated_at,
        )

    @router.get("", response_model=List[BoardResponse], status_code=status.HTTP_200_OK)
    def list_boards():
        res = list_boards_uc.exec(ListBoardsInput())
        return [
            BoardResponse(
                id=b.id,
                name=b.name,
                description=b.description,
                columns=[
                    ColumnResponse(
                        id=c.id,
                        name=c.name,
                        order_index=c.order_index,
                        wip_limit=c.wip_limit,
                    )
                    for c in b.columns
                ],
                created_at=b.created_at,
                updated_at=b.updated_at,
            )
            for b in res.boards
        ]

    @router.get("/{board_id}", response_model=BoardResponse, status_code=status.HTTP_200_OK)
    def get_board(board_id: str):
        res = get_board_uc.exec(GetBoardInput(board_id=board_id))
        return BoardResponse(
            id=res.id,
            name=res.name,
            description=res.description,
            columns=[
                ColumnResponse(
                    id=c.id,
                    name=c.name,
                    order_index=c.order_index,
                    wip_limit=c.wip_limit,
                )
                for c in res.columns
            ],
            created_at=res.created_at,
            updated_at=res.updated_at,
        )

    @router.post("/{board_id}/columns", response_model=ColumnResponse, status_code=status.HTTP_201_CREATED)
    def add_column(board_id: str, payload: AddColumnRequest):
        res = add_column_uc.exec(
            AddColumnInput(
                board_id=board_id,
                name=payload.name,
                wip_limit=payload.wip_limit,
            )
        )
        return ColumnResponse(
            id=res.column.id,
            name=res.column.name,
            order_index=res.column.order_index,
            wip_limit=res.column.wip_limit,
        )

    @router.post("/{board_id}/tasks", response_model=TaskResponse, status_code=status.HTTP_201_CREATED)
    def create_task(board_id: str, payload: CreateTaskRequest):
        res = create_task_uc.exec(
            CreateTaskInput(
                board_id=board_id,
                column=payload.column,
                title=payload.title,
                description=payload.description,
                priority=payload.priority,
                due_date=payload.due_date,
            )
        )
        return TaskResponse(
            id=res.id,
            board_id=res.board_id,
            column_id=res.column_id,
            column_name=res.column_name,
            title=res.title,
            description=res.description,
            priority=res.priority,
            due_date=res.due_date,
            created_at=res.created_at,
            updated_at=res.updated_at,
        )

    @router.get("/{board_id}/tasks", response_model=List[TaskResponse], status_code=status.HTTP_200_OK)
    def list_tasks(board_id: str, column_id: Optional[str] = Query(None)):
        res = list_tasks_uc.exec(ListTasksInput(board_id=board_id, column_id=column_id))
        return [
            TaskResponse(
                id=t.id,
                board_id=t.board_id,
                column_id=t.column_id,
                column_name=t.column_name,
                title=t.title,
                description=t.description,
                priority=t.priority,
                due_date=t.due_date,
                created_at=t.created_at,
                updated_at=t.updated_at,
            )
            for t in res.tasks
        ]

    return router

