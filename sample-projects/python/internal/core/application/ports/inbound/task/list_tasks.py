from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


@dataclass(frozen=True)
class TaskItemDTO:
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


@dataclass(frozen=True)
class ListTasksInput:
    board_id: str
    column_id: Optional[str] = None
    column_name: Optional[str] = None


@dataclass(frozen=True)
class ListTasksOutput:
    board_id: str
    tasks: List[TaskItemDTO] = field(default_factory=list)


class InboundListTasksBoardNotFoundError(Exception):
    def __init__(self, board_id: str):
        super().__init__(f"Board with id '{board_id}' was not found.")
        self.code = "BOARD_NOT_FOUND"
        self.board_id = board_id


class ListTasksUseCase(ABC):
    @abstractmethod
    def exec(self, req: ListTasksInput) -> ListTasksOutput:
        pass

