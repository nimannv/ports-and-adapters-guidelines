import pytest
from internal.core.domain.board.entities import Board, Column
from internal.core.domain.board.errors import (
    ColumnNotFoundError,
    DuplicateColumnNameError,
    InvalidBoardError,
    InvalidWipLimitError,
)
from internal.core.domain.task.entities import Priority, Task
from internal.core.domain.task.errors import (
    InvalidTaskError,
    InvalidTaskMovementError,
)


def test_create_board_defaults():
    board = Board.create(name="Sprint 1", description="Release MVP")
    assert board.name == "Sprint 1"
    assert board.description == "Release MVP"
    assert len(board.columns) == 3
    col_names = [c.name for c in board.columns]
    assert col_names == ["Backlog", "In Progress", "Done"]

    in_progress = board.columns[1]
    assert in_progress.wip_limit == 3
    assert in_progress.is_full(2) is False
    assert in_progress.is_full(3) is True


def test_board_name_empty_fails():
    with pytest.raises(InvalidBoardError):
        Board.create(name="   ")


def test_column_invalid_wip_limit():
    with pytest.raises(InvalidWipLimitError):
        Column(id="col-1", name="Review", order_index=0, wip_limit=0)

    with pytest.raises(InvalidWipLimitError):
        Column(id="col-2", name="Review", order_index=0, wip_limit=-1)


def test_board_duplicate_column_name():
    board = Board.create(name="Sprint 1")
    with pytest.raises(DuplicateColumnNameError):
        board.add_column("Backlog")

    with pytest.raises(DuplicateColumnNameError):
        board.add_column("backlog")


def test_board_add_column_success():
    board = Board.create(name="Sprint 1")
    col = board.add_column("QA Review", wip_limit=2)
    assert col.name == "QA Review"
    assert col.wip_limit == 2
    assert col.order_index == 3
    assert board.has_column_id(col.id) is True


def test_board_get_column():
    board = Board.create(name="Sprint 1")
    col = board.get_column("Backlog")
    assert col.name == "Backlog"

    col_by_id = board.get_column(col.id)
    assert col_by_id.id == col.id

    with pytest.raises(ColumnNotFoundError):
        board.get_column("NonExistent")


def test_task_creation_and_movement():
    task = Task.create(
        board_id="b-1",
        column_id="c-1",
        title="Implement Login",
        priority=Priority.HIGH,
    )
    assert task.title == "Implement Login"
    assert task.priority == Priority.HIGH
    assert task.column_id == "c-1"

    # Move to column on the same board
    task.move_to(new_column_id="c-2", target_board_id="b-1")
    assert task.column_id == "c-2"

    # Attempt move across boards fails
    with pytest.raises(InvalidTaskMovementError):
        task.move_to(new_column_id="c-3", target_board_id="different-board")


def test_task_empty_title_fails():
    with pytest.raises(InvalidTaskError):
        Task.create(board_id="b-1", column_id="c-1", title="  ")

