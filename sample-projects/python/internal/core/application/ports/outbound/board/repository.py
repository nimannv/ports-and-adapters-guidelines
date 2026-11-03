from abc import ABC, abstractmethod
from typing import List, Optional

from internal.core.application.ports.outbound.board.dtos import BoardRecord, ColumnRecord


class BoardRepository(ABC):
    """Outbound port for persisting and retrieving board records."""

    @abstractmethod
    def save(self, record: BoardRecord) -> BoardRecord:
        """Persist a new or updated board with its columns."""
        pass

    @abstractmethod
    def get_by_id(self, board_id: str) -> Optional[BoardRecord]:
        """Retrieve a board by its UUID string with its columns."""
        pass

    @abstractmethod
    def list_all(self) -> List[BoardRecord]:
        """Retrieve all boards."""
        pass

    @abstractmethod
    def add_column(self, board_id: str, column: ColumnRecord) -> ColumnRecord:
        """Append a new column to an existing board."""
        pass

