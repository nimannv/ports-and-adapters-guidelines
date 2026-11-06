from dataclasses import dataclass, field
from datetime import datetime, timezone
from enum import Enum
from typing import Optional
import uuid

from internal.core.domain.task.errors import InvalidTaskError, InvalidTaskMovementError


class Priority(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"
    URGENT = "URGENT"

    @classmethod
    def from_str(cls, val: str) -> "Priority":
        try:
            return cls[val.upper()]
        except KeyError:
            valid = ", ".join([e.value for e in cls])
            raise InvalidTaskError(f"Invalid priority '{val}'. Must be one of: {valid}")


@dataclass
class Task:
    id: str
    board_id: str
    column_id: str
    title: str
    description: str = ""
    priority: Priority = Priority.MEDIUM
    due_date: Optional[datetime] = None
    created_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = field(default_factory=lambda: datetime.now(timezone.utc))

    def __post_init__(self) -> None:
        if not self.title or not self.title.strip():
            raise InvalidTaskError("Task title cannot be empty.")
        if not self.board_id or not self.board_id.strip():
            raise InvalidTaskError("Task must be associated with a valid board.")
        if not self.column_id or not self.column_id.strip():
            raise InvalidTaskError("Task must be associated with a valid column.")

    @classmethod
    def create(
        cls,
        board_id: str,
        column_id: str,
        title: str,
        description: str = "",
        priority: Priority = Priority.MEDIUM,
        due_date: Optional[datetime] = None,
    ) -> "Task":
        now = datetime.now(timezone.utc)
        return cls(
            id=str(uuid.uuid4()),
            board_id=board_id.strip(),
            column_id=column_id.strip(),
            title=title.strip(),
            description=description.strip(),
            priority=priority,
            due_date=due_date,
            created_at=now,
            updated_at=now,
        )

    def move_to(self, new_column_id: str, target_board_id: str) -> None:
        if target_board_id != self.board_id:
            raise InvalidTaskMovementError("Tasks can only be moved between columns that belong to the same board.")
        self.column_id = new_column_id
        self.updated_at = datetime.now(timezone.utc)

