from datetime import datetime
from typing import List, Optional
from sqlalchemy.orm import sessionmaker

from internal.adapters.outbound.sql.mappers import task_model_to_record
from internal.adapters.outbound.sql.models import TaskModel
from internal.core.application.ports.outbound.task import (
    OutboundTaskRepositoryError,
    TaskRecord,
    TaskRepository,
)


class SqlAlchemyTaskRepository(TaskRepository):
    def __init__(self, session_factory: sessionmaker):
        self._session_factory = session_factory

    def save(self, record: TaskRecord) -> TaskRecord:
        try:
            with self._session_factory() as session:
                task_model = TaskModel(
                    id=record.id,
                    board_id=record.board_id,
                    column_id=record.column_id,
                    title=record.title,
                    description=record.description,
                    priority=record.priority,
                    due_date=record.due_date,
                    created_at=record.created_at,
                    updated_at=record.updated_at,
                )
                session.add(task_model)
                session.commit()
                return task_model_to_record(task_model)
        except Exception as e:
            raise OutboundTaskRepositoryError(f"Database error while saving task: {e}") from e

    def get_by_id(self, task_id: str) -> Optional[TaskRecord]:
        try:
            with self._session_factory() as session:
                task = session.query(TaskModel).filter(TaskModel.id == task_id).first()
                if not task:
                    return None
                return task_model_to_record(task)
        except Exception as e:
            raise OutboundTaskRepositoryError(f"Database error retrieving task '{task_id}': {e}") from e

    def list_by_board(self, board_id: str, column_id: Optional[str] = None) -> List[TaskRecord]:
        try:
            with self._session_factory() as session:
                query = session.query(TaskModel).filter(TaskModel.board_id == board_id)
                if column_id:
                    query = query.filter(TaskModel.column_id == column_id)
                tasks = query.order_by(TaskModel.created_at.asc()).all()
                return [task_model_to_record(t) for t in tasks]
        except Exception as e:
            raise OutboundTaskRepositoryError(f"Database error listing tasks for board '{board_id}': {e}") from e

    def count_by_column(self, column_id: str) -> int:
        try:
            with self._session_factory() as session:
                return session.query(TaskModel).filter(TaskModel.column_id == column_id).count()
        except Exception as e:
            raise OutboundTaskRepositoryError(f"Database error counting tasks for column '{column_id}': {e}") from e

    def update_column(self, task_id: str, new_column_id: str, updated_at: datetime) -> Optional[TaskRecord]:
        try:
            with self._session_factory() as session:
                task = session.query(TaskModel).filter(TaskModel.id == task_id).first()
                if not task:
                    return None
                task.column_id = new_column_id
                task.updated_at = updated_at
                session.commit()
                session.refresh(task)
                return task_model_to_record(task)
        except Exception as e:
            raise OutboundTaskRepositoryError(f"Database error updating column for task '{task_id}': {e}") from e

    def delete(self, task_id: str) -> bool:
        try:
            with self._session_factory() as session:
                task = session.query(TaskModel).filter(TaskModel.id == task_id).first()
                if not task:
                    return False
                session.delete(task)
                session.commit()
                return True
        except Exception as e:
            raise OutboundTaskRepositoryError(f"Database error deleting task '{task_id}': {e}") from e

