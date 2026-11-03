from internal.core.application.ports.outbound.board.dtos import BoardRecord, ColumnRecord
from internal.core.application.ports.outbound.task.dtos import TaskRecord
from internal.adapters.outbound.sql.models import BoardModel, ColumnModel, TaskModel


def column_model_to_record(model: ColumnModel) -> ColumnRecord:
    return ColumnRecord(
        id=model.id,
        board_id=model.board_id,
        name=model.name,
        order_index=model.order_index,
        wip_limit=model.wip_limit,
    )


def board_model_to_record(model: BoardModel) -> BoardRecord:
    return BoardRecord(
        id=model.id,
        name=model.name,
        description=model.description,
        columns=[column_model_to_record(col) for col in model.columns],
        created_at=model.created_at,
        updated_at=model.updated_at,
    )


def task_model_to_record(model: TaskModel) -> TaskRecord:
    return TaskRecord(
        id=model.id,
        board_id=model.board_id,
        column_id=model.column_id,
        title=model.title,
        description=model.description,
        priority=model.priority,
        due_date=model.due_date,
        created_at=model.created_at,
        updated_at=model.updated_at,
    )

