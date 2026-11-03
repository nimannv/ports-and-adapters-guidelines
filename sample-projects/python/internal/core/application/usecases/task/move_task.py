from datetime import datetime, timezone

from internal.core.application.ports.inbound.task.move_task import (
    InboundInvalidTaskMovementError,
    InboundMoveColumnNotFoundError,
    InboundMoveWipLimitExceededError,
    InboundTaskNotFoundError,
    MoveTaskInput,
    MoveTaskOutput,
    MoveTaskUseCase,
)
from internal.core.application.ports.outbound.board import BoardRepository
from internal.core.application.ports.outbound.task import TaskRepository


class MoveTaskService(MoveTaskUseCase):
    def __init__(self, board_repo: BoardRepository, task_repo: TaskRepository):
        self._board_repo = board_repo
        self._task_repo = task_repo

    def exec(self, req: MoveTaskInput) -> MoveTaskOutput:
        task = self._task_repo.get_by_id(req.task_id)
        if not task:
            raise InboundTaskNotFoundError(req.task_id)

        board = self._board_repo.get_by_id(task.board_id)
        if not board:
            raise InboundInvalidTaskMovementError(f"Parent board '{task.board_id}' no longer exists.")

        # Resolve target column
        target_column = None
        target_ref = req.target_column_id or req.target_column_name
        if not target_ref:
            raise InboundInvalidTaskMovementError("Either target_column_id or target_column_name must be specified.")

        target_ref_lower = target_ref.strip().lower()
        for col in board.columns:
            if col.id == target_ref or col.name.strip().lower() == target_ref_lower:
                target_column = col
                break

        if not target_column:
            raise InboundMoveColumnNotFoundError(target_ref, board.id)

        # Check if moving to the same column
        if task.column_id == target_column.id:
            return MoveTaskOutput(
                task_id=task.id,
                board_id=task.board_id,
                column_id=task.column_id,
                column_name=target_column.name,
                title=task.title,
                priority=task.priority,
                updated_at=task.updated_at,
            )

        # Check WIP limit of destination column
        if target_column.wip_limit is not None:
            current_count = self._task_repo.count_by_column(target_column.id)
            if current_count >= target_column.wip_limit:
                raise InboundMoveWipLimitExceededError(
                    message=f"Cannot move task to '{target_column.name}': column limit of {target_column.wip_limit} tasks has been reached.",
                    column_id=target_column.id,
                    current_count=current_count,
                    limit=target_column.wip_limit,
                )

        now = datetime.now(timezone.utc)
        updated_task = self._task_repo.update_column(
            task_id=task.id,
            new_column_id=target_column.id,
            updated_at=now,
        )

        if not updated_task:
            raise InboundTaskNotFoundError(task.id)

        return MoveTaskOutput(
            task_id=updated_task.id,
            board_id=updated_task.board_id,
            column_id=updated_task.column_id,
            column_name=target_column.name,
            title=updated_task.title,
            priority=updated_task.priority,
            updated_at=updated_task.updated_at,
        )

