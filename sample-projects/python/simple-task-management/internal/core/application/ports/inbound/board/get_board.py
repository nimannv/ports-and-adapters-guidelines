from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional

from internal.core.application.ports.inbound.board.create_board import ColumnDTO


@dataclass(frozen=True)
class GetBoardInput:
    board_id: str


@dataclass(frozen=True)
class GetBoardOutput:
    id: str
    name: str
    description: str
    columns: List[ColumnDTO] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class BoardNotFoundError(Exception):
    def __init__(self, board_id: str):
        super().__init__(f"Board with id '{board_id}' was not found.")
        self.code = "BOARD_NOT_FOUND"
        self.board_id = board_id


class GetBoardUseCase(ABC):
    @abstractmethod
    def exec(self, req: GetBoardInput) -> GetBoardOutput:
        pass

