from dataclasses import dataclass, field
from datetime import datetime
from typing import List, Optional


@dataclass(frozen=True)
class ColumnRecord:
    id: str
    board_id: str
    name: str
    order_index: int
    wip_limit: Optional[int] = None


@dataclass(frozen=True)
class BoardRecord:
    id: str
    name: str
    description: str
    columns: List[ColumnRecord] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None

