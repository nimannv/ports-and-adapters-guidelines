from internal.core.application.ports.inbound.task.list_tasks import (
    InboundListTasksBoardNotFoundError,
    ListTasksInput,
    ListTasksOutput,
    ListTasksUseCase,
    TaskItemDTO,
)
from internal.core.application.ports.outbound.board import BoardRepository
from internal.core.application.ports.outbound.task import TaskRepository


class ListTasksService(ListTasksUseCase):
    def __init__(self, board_repo: BoardRepository, task_repo: TaskRepository):
        self._board_repo = board_repo
        self._task_repo = task_repo

    def exec(self, req: ListTasksInput) -> ListTasksOutput:
        board = self._board_repo.get_by_id(req.board_id)
        if not board:
            raise InboundListTasksBoardNotFoundError(req.board_id)

        col_map = {col.id: col.name for col in board.columns}
        col_name_to_id = {col.name.strip().lower(): col.id for col in board.columns}

        filter_column_id = req.column_id
        if not filter_column_id and req.column_name:
            filter_column_id = col_name_to_id.get(req.column_name.strip().lower())

        records = self._task_repo.list_by_board(req.board_id, column_id=filter_column_id)

        task_dtos = [
            TaskItemDTO(
                id=t.id,
                board_id=t.board_id,
                column_id=t.column_id,
                column_name=col_map.get(t.column_id, "Unknown"),
                title=t.title,
                description=t.description,
                priority=t.priority,
                due_date=t.due_date,
                created_at=t.created_at,
                updated_at=t.updated_at,
            )
            for t in records
        ]

        return ListTasksOutput(board_id=req.board_id, tasks=task_dtos)

