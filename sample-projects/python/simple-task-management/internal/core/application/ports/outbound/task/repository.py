from abc import ABC, abstractmethod
from datetime import datetime
from typing import List, Optional

from internal.core.application.ports.outbound.task.dtos import TaskRecord


class TaskRepository(ABC):
    """Outbound port for persisting, querying, and updating tasks."""

    @abstractmethod
    def save(self, record: TaskRecord) -> TaskRecord:
        """Persist a new task."""
        pass

    @abstractmethod
    def get_by_id(self, task_id: str) -> Optional[TaskRecord]:
        """Fetch a task by id."""
        pass

    @abstractmethod
    def list_by_board(self, board_id: str, column_id: Optional[str] = None) -> List[TaskRecord]:
        """List tasks for a board, optionally filtered by column."""
        pass

    @abstractmethod
    def count_by_column(self, column_id: str) -> int:
        """Count the active tasks assigned to a specific column."""
        pass

    @abstractmethod
    def update_column(self, task_id: str, new_column_id: str, updated_at: datetime) -> Optional[TaskRecord]:
        """Update a task's destination column and timestamp."""
        pass

    @abstractmethod
    def delete(self, task_id: str) -> bool:
        """Delete a task by id. Returns True if deleted, False if not found."""
        pass

