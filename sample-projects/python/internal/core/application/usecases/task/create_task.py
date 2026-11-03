from datetime import datetime, timezone

from internal.core.application.ports.inbound.task.create_task import (
    CreateTaskInput,
    CreateTaskOutput,
    CreateTaskUseCase,
    InboundBoardNotFoundError,
    InboundColumnNotFoundError,
    InboundInvalidTaskError,
    InboundWipLimitExceededError,
)
from internal.core.application.ports.outbound.board import BoardRepository
from internal.core.application.ports.outbound.task import TaskRecord, TaskRepository
from internal.core.domain.task.entities import Priority, Task
from internal.core.domain.task.errors import InvalidTaskError


class CreateTaskService(CreateTaskUseCase):
    def __init__(self, board_repo: BoardRepository, task_repo: TaskRepository):
        self._board_repo = board_repo
        self._task_repo = task_repo

    def exec(self, req: CreateTaskInput) -> CreateTaskOutput:
        board = self._board_repo.get_by_id(req.board_id)
        if not board:
            raise InboundBoardNotFoundError(req.board_id)

        # Resolve column by ID or Name
        target_column = None
        req_col_lower = req.column.strip().lower()
        for col in board.columns:
            if col.id == req.column or col.name.strip().lower() == req_col_lower:
                target_column = col
                break

        if not target_column:
            raise InboundColumnNotFoundError(req.column, req.board_id)

        # Check WIP Limit
        if target_column.wip_limit is not None:
            current_count = self._task_repo.count_by_column(target_column.id)
            if current_count >= target_column.wip_limit:
                raise InboundWipLimitExceededError(
                    message=f"Cannot add task to '{target_column.name}': column limit of {target_column.wip_limit} tasks has been reached.",
                    column_id=target_column.id,
                    current_count=current_count,
                    limit=target_column.wip_limit,
                )

        try:
            priority = Priority.from_str(req.priority)
            domain_task = Task.create(
                board_id=board.id,
                column_id=target_column.id,
                title=req.title,
                description=req.description,
                priority=priority,
                due_date=req.due_date,
            )
        except InvalidTaskError as e:
            raise InboundInvalidTaskError(e.message) from e

        record = TaskRecord(
            id=domain_task.id,
            board_id=domain_task.board_id,
            column_id=domain_task.column_id,
            title=domain_task.title,
            description=domain_task.description,
            priority=domain_task.priority.value,
            due_date=domain_task.due_date,
            created_at=domain_task.created_at,
            updated_at=domain_task.updated_at,
        )

        saved = self._task_repo.save(record)

        return CreateTaskOutput(
            id=saved.id,
            board_id=saved.board_id,
            column_id=saved.column_id,
            column_name=target_column.name,
            title=saved.title,
            description=saved.description,
            priority=saved.priority,
            due_date=saved.due_date,
            created_at=saved.created_at,
            updated_at=saved.updated_at,
        )

