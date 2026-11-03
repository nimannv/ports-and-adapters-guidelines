from abc import ABC, abstractmethod
from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass(frozen=True)
class MoveTaskInput:
    task_id: str
    target_column_id: Optional[str] = None
    target_column_name: Optional[str] = None


@dataclass(frozen=True)
class MoveTaskOutput:
    task_id: str
    board_id: str
    column_id: str
    column_name: str
    title: str
    priority: str
    updated_at: Optional[datetime] = None


class InboundMoveTaskError(Exception):
    def __init__(self, message: str, code: str = "MOVE_TASK_ERROR"):
        super().__init__(message)
        self.message = message
        self.code = code


class InboundTaskNotFoundError(InboundMoveTaskError):
    def __init__(self, task_id: str):
        super().__init__(f"Task with id '{task_id}' was not found.", code="TASK_NOT_FOUND")
        self.task_id = task_id


class InboundMoveColumnNotFoundError(InboundMoveTaskError):
    def __init__(self, column_ref: str, board_id: Optional[str] = None):
        msg = f"Target column '{column_ref}' was not found"
        if board_id:
            msg += f" on board '{board_id}'."
        else:
            msg += "."
        super().__init__(msg, code="COLUMN_NOT_FOUND")
        self.column_ref = column_ref
        self.board_id = board_id


class InboundMoveWipLimitExceededError(InboundMoveTaskError):
    def __init__(self, message: str, column_id: str, current_count: int, limit: int):
        super().__init__(message, code="WIP_LIMIT_EXCEEDED")
        self.column_id = column_id
        self.current_count = current_count
        self.limit = limit


class InboundInvalidTaskMovementError(InboundMoveTaskError):
    def __init__(self, message: str):
        super().__init__(message, code="INVALID_TASK_MOVEMENT")


class MoveTaskUseCase(ABC):
    @abstractmethod
    def exec(self, req: MoveTaskInput) -> MoveTaskOutput:
        pass

