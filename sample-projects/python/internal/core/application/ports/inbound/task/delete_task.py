from abc import ABC, abstractmethod
from dataclasses import dataclass


@dataclass(frozen=True)
class DeleteTaskInput:
    task_id: str


@dataclass(frozen=True)
class DeleteTaskOutput:
    task_id: str
    deleted: bool


class InboundDeleteTaskNotFoundError(Exception):
    def __init__(self, task_id: str):
        super().__init__(f"Task with id '{task_id}' was not found.")
        self.code = "TASK_NOT_FOUND"
        self.task_id = task_id


class DeleteTaskUseCase(ABC):
    @abstractmethod
    def exec(self, req: DeleteTaskInput) -> DeleteTaskOutput:
        pass

