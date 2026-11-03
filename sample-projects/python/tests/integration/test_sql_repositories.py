import uuid
from internal.adapters.outbound.sql import (
    SqlAlchemyBoardRepository,
    SqlAlchemyTaskRepository,
    create_db_engine,
    create_session_factory,
    init_db,
)
from internal.core.application.ports.outbound.board import BoardRecord, ColumnRecord
from internal.core.application.ports.outbound.task import TaskRecord


def test_sql_repositories():
    engine = create_db_engine("sqlite:///:memory:")
    init_db(engine)
    session_factory = create_session_factory(engine)

    board_repo = SqlAlchemyBoardRepository(session_factory)
    task_repo = SqlAlchemyTaskRepository(session_factory)

    # 1. Save Board with columns
    board_id = str(uuid.uuid4())
    col1_id = str(uuid.uuid4())
    col2_id = str(uuid.uuid4())

    board = BoardRecord(
        id=board_id,
        name="Test Board",
        description="Testing SQL repo",
        columns=[
            ColumnRecord(id=col1_id, board_id=board_id, name="To Do", order_index=0, wip_limit=None),
            ColumnRecord(id=col2_id, board_id=board_id, name="Doing", order_index=1, wip_limit=2),
        ],
    )
    saved_board = board_repo.save(board)
    assert saved_board.id == board_id
    assert len(saved_board.columns) == 2

    # 2. Get board by ID
    found_board = board_repo.get_by_id(board_id)
    assert found_board is not None
    assert found_board.name == "Test Board"

    # 3. Add column to board
    col3_id = str(uuid.uuid4())
    added_col = board_repo.add_column(
        board_id,
        ColumnRecord(id=col3_id, board_id=board_id, name="Done", order_index=2, wip_limit=None),
    )
    assert added_col.name == "Done"

    refreshed_board = board_repo.get_by_id(board_id)
    assert len(refreshed_board.columns) == 3

    # 4. Save Task
    task_id = str(uuid.uuid4())
    task = TaskRecord(
        id=task_id,
        board_id=board_id,
        column_id=col1_id,
        title="Test Task",
        description="Testing task storage",
        priority="HIGH",
    )
    saved_task = task_repo.save(task)
    assert saved_task.id == task_id
    assert task_repo.count_by_column(col1_id) == 1

    # 5. Update Task column
    from datetime import datetime, timezone

    updated = task_repo.update_column(task_id, col2_id, datetime.now(timezone.utc))
    assert updated.column_id == col2_id
    assert task_repo.count_by_column(col1_id) == 0
    assert task_repo.count_by_column(col2_id) == 1

    # 6. List tasks
    tasks = task_repo.list_by_board(board_id)
    assert len(tasks) == 1

    # 7. Delete task
    assert task_repo.delete(task_id) is True
    assert task_repo.get_by_id(task_id) is None

