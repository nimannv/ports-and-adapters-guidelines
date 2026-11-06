from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


@dataclass(frozen=True)
class ColumnDTO:
    id: str
    name: str
    order_index: int
    wip_limit: Optional[int] = None


@dataclass(frozen=True)
class CreateBoardInput:
    name: str
    description: str = ""
    custom_columns: Optional[List[dict]] = None


@dataclass(frozen=True)
class CreateBoardOutput:
    id: str
    name: str
    description: str
    columns: List[ColumnDTO] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class CreateBoardError(Exception):
    def __init__(self, message: str, code: str = "CREATE_BOARD_ERROR"):
        super().__init__(message)
        self.message = message
        self.code = code


class CreateBoardUseCase(ABC):
    @abstractmethod
    def exec(self, req: CreateBoardInput) -> CreateBoardOutput:
        pass

