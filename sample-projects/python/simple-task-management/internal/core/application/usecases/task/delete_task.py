from internal.core.application.ports.inbound.task.delete_task import (
    DeleteTaskInput,
    DeleteTaskOutput,
    DeleteTaskUseCase,
    InboundDeleteTaskNotFoundError,
)
from internal.core.application.ports.outbound.task import TaskRepository


class DeleteTaskService(DeleteTaskUseCase):
    def __init__(self, task_repo: TaskRepository):
        self._task_repo = task_repo

    def exec(self, req: DeleteTaskInput) -> DeleteTaskOutput:
        existing = self._task_repo.get_by_id(req.task_id)
        if not existing:
            raise InboundDeleteTaskNotFoundError(req.task_id)

        deleted = self._task_repo.delete(req.task_id)
        return DeleteTaskOutput(task_id=req.task_id, deleted=deleted)

