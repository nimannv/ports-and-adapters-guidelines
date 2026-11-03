from abc import ABC, abstractmethod
from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass(frozen=True)
class CreateTaskInput:
    board_id: str
    column: str  # Can be column ID or column name
    title: str
    description: str = ""
    priority: str = "MEDIUM"
    due_date: Optional[datetime] = None


@dataclass(frozen=True)
class CreateTaskOutput:
    id: str
    board_id: str
    column_id: str
    column_name: str
    title: str
    description: str
    priority: str
    due_date: Optional[datetime] = None
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class InboundTaskError(Exception):
    def __init__(self, message: str, code: str = "TASK_ERROR"):
        super().__init__(message)
        self.message = message
        self.code = code


class InboundBoardNotFoundError(InboundTaskError):
    def __init__(self, board_id: str):
        super().__init__(f"Board with id '{board_id}' was not found.", code="BOARD_NOT_FOUND")
        self.board_id = board_id


class InboundColumnNotFoundError(InboundTaskError):
    def __init__(self, column_ref: str, board_id: Optional[str] = None):
        msg = f"Column '{column_ref}' was not found"
        if board_id:
            msg += f" on board '{board_id}'."
        else:
            msg += "."
        super().__init__(msg, code="COLUMN_NOT_FOUND")
        self.column_ref = column_ref
        self.board_id = board_id


class InboundWipLimitExceededError(InboundTaskError):
    def __init__(self, message: str, column_id: str, current_count: int, limit: int):
        super().__init__(message, code="WIP_LIMIT_EXCEEDED")
        self.column_id = column_id
        self.current_count = current_count
        self.limit = limit


class InboundInvalidTaskError(InboundTaskError):
    def __init__(self, message: str):
        super().__init__(message, code="INVALID_TASK")


class CreateTaskUseCase(ABC):
    @abstractmethod
    def exec(self, req: CreateTaskInput) -> CreateTaskOutput:
        pass

