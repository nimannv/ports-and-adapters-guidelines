from fastapi import APIRouter, Response, status

from internal.adapters.inbound.http.dtos import MoveTaskRequest, MoveTaskResponse
from internal.core.application.ports.inbound.task import (
    DeleteTaskInput,
    DeleteTaskUseCase,
    MoveTaskInput,
    MoveTaskUseCase,
)


def create_task_router(
    move_task_uc: MoveTaskUseCase,
    delete_task_uc: DeleteTaskUseCase,
) -> APIRouter:
    router = APIRouter(prefix="/api/v1/tasks", tags=["tasks"])

    @router.patch("/{task_id}/move", response_model=MoveTaskResponse, status_code=status.HTTP_200_OK)
    def move_task(task_id: str, payload: MoveTaskRequest):
        res = move_task_uc.exec(
            MoveTaskInput(
                task_id=task_id,
                target_column_id=payload.target_column_id,
                target_column_name=payload.target_column_name,
            )
        )
        return MoveTaskResponse(
            task_id=res.task_id,
            board_id=res.board_id,
            column_id=res.column_id,
            column_name=res.column_name,
            title=res.title,
            priority=res.priority,
            updated_at=res.updated_at,
        )

    @router.delete("/{task_id}", status_code=status.HTTP_204_NO_CONTENT)
    def delete_task(task_id: str):
        delete_task_uc.exec(DeleteTaskInput(task_id=task_id))
        return Response(status_code=status.HTTP_204_NO_CONTENT)

    return router

