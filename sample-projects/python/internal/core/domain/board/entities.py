from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import List, Optional
import uuid

from internal.core.domain.board.errors import (
    DuplicateColumnNameError,
    InvalidBoardError,
    InvalidWipLimitError,
    ColumnNotFoundError,
)


@dataclass
class Column:
    id: str
    name: str
    order_index: int
    wip_limit: Optional[int] = None

    def __post_init__(self) -> None:
        if not self.name or not self.name.strip():
            raise InvalidBoardError("Column name cannot be empty.")
        if self.wip_limit is not None and self.wip_limit <= 0:
            raise InvalidWipLimitError(self.wip_limit)

    def is_full(self, current_task_count: int) -> bool:
        if self.wip_limit is None:
            return False
        return current_task_count >= self.wip_limit


@dataclass
class Board:
    id: str
    name: str
    description: str
    columns: List[Column] = field(default_factory=list)
    created_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))

    def __post_init__(self) -> None:
        if not self.name or not self.name.strip():
            raise InvalidBoardError("Board name cannot be empty.")
        self._validate_unique_columns(self.columns)

    def _validate_unique_columns(self, cols: List[Column]) -> None:
        seen = set()
        for col in cols:
            normalized = col.name.strip().lower()
            if normalized in seen:
                raise DuplicateColumnNameError(col.name, self.id)
            seen.add(normalized)

    @classmethod
    def create(
        cls,
        name: str,
        description: str = "",
        custom_columns: Optional[List[dict]] = None,
    ) -> "Board":
        board_id = str(uuid.uuid4())
        now = datetime.now(timezone.utc)

        if custom_columns:
            columns: List[Column] = []
            for idx, col_data in enumerate(custom_columns):
                columns.append(
                    Column(
                        id=str(uuid.uuid4()),
                        name=col_data["name"],
                        order_index=col_data.get("order_index", idx),
                        wip_limit=col_data.get("wip_limit"),
                    )
                )
        else:
            # Default columns per specification: Backlog (unlimited), In Progress (limit: 3), Done (unlimited)
            columns = [
                Column(id=str(uuid.uuid4()), name="Backlog", order_index=0, wip_limit=None),
                Column(id=str(uuid.uuid4()), name="In Progress", order_index=1, wip_limit=3),
                Column(id=str(uuid.uuid4()), name="Done", order_index=2, wip_limit=None),
            ]

        if not columns:
            raise InvalidBoardError("A board must contain at least one column upon creation.")

        board = cls(
            id=board_id,
            name=name.strip(),
            description=description.strip(),
            columns=columns,
            created_at=now,
            updated_at=now,
        )
        return board

    def add_column(self, name: str, wip_limit: Optional[int] = None) -> Column:
        cleaned_name = name.strip()
        if not cleaned_name:
            raise InvalidBoardError("Column name cannot be empty.")

        for col in self.columns:
            if col.name.strip().lower() == cleaned_name.lower():
                raise DuplicateColumnNameError(cleaned_name, self.id)

        next_order = len(self.columns)
        new_column = Column(
            id=str(uuid.uuid4()),
            name=cleaned_name,
            order_index=next_order,
            wip_limit=wip_limit,
        )
        self.columns.append(new_column)
        self.updated_at = datetime.now(timezone.utc)
        return new_column

    def get_column(self, column_ref: str) -> Column:
        """Find column by either ID or case-insensitive Name."""
        ref_lower = column_ref.strip().lower()
        for col in self.columns:
            if col.id == column_ref or col.name.strip().lower() == ref_lower:
                return col
        raise ColumnNotFoundError(column_ref, self.id)

    def has_column_id(self, column_id: str) -> bool:
        return any(col.id == column_id for col in self.columns)

