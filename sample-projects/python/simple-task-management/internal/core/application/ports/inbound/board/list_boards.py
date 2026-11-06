from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional

from internal.core.application.ports.inbound.board.create_board import ColumnDTO


@dataclass(frozen=True)
class BoardSummaryDTO:
    id: str
    name: str
    description: str
    columns: List[ColumnDTO] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass(frozen=True)
class ListBoardsInput:
    pass


@dataclass(frozen=True)
class ListBoardsOutput:
    boards: List[BoardSummaryDTO] = field(default_factory=list)


class ListBoardsUseCase(ABC):
    @abstractmethod
    def exec(self, req: ListBoardsInput) -> ListBoardsOutput:
        pass

