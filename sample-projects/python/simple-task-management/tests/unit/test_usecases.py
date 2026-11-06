from datetime import datetime, timezone
from typing import Dict, List, Optional
import pytest

from internal.core.application.ports.inbound.board import (
    AddColumnInput,
    BoardNotFoundError,
    CreateBoardInput,
    GetBoardInput,
    InboundDuplicateColumnNameError,
    ListBoardsInput,
)
from internal.core.application.ports.inbound.task import (
    CreateTaskInput,
    DeleteTaskInput,
    InboundBoardNotFoundError,
    InboundColumnNotFoundError,
    InboundDeleteTaskNotFoundError,
    InboundMoveWipLimitExceededError,
    InboundTaskNotFoundError,
    ListTasksInput,
    MoveTaskInput,
)
from internal.core.application.ports.outbound.board import BoardRecord, BoardRepository, ColumnRecord
from internal.core.application.ports.outbound.task import TaskRecord, TaskRepository
from internal.core.application.usecases.board import (
    AddColumnService,
    CreateBoardService,
    GetBoardService,
    ListBoardsService,
)
from internal.core.application.usecases.task import (
    CreateTaskService,
    DeleteTaskService,
    ListTasksService,
    MoveTaskService,
)


class InMemoryBoardRepository(BoardRepository):
    def __init__(self):
        self.boards: Dict[str, BoardRecord] = {}

    def save(self, record: BoardRecord) -> BoardRecord:
        self.boards[record.id] = record
        return record

    def get_by_id(self, board_id: str) -> Optional[BoardRecord]:
        return self.boards.get(board_id)

    def list_all(self) -> List[BoardRecord]:
        return list(self.boards.values())

    def add_column(self, board_id: str, column: ColumnRecord) -> ColumnRecord:
        board = self.boards[board_id]
        updated_cols = list(board.columns) + [column]
        self.boards[board_id] = BoardRecord(
            id=board.id,
            name=board.name,
            description=board.description,
            columns=updated_cols,
            created_at=board.created_at,
            updated_at=datetime.now(timezone.utc),
        )
        return column


class InMemoryTaskRepository(TaskRepository):
    def __init__(self):
        self.tasks: Dict[str, TaskRecord] = {}

    def save(self, record: TaskRecord) -> TaskRecord:
        self.tasks[record.id] = record
        return record

    def get_by_id(self, task_id: str) -> Optional[TaskRecord]:
        return self.tasks.get(task_id)

    def list_by_board(self, board_id: str, column_id: Optional[str] = None) -> List[TaskRecord]:
        results = [t for t in self.tasks.values() if t.board_id == board_id]
        if column_id:
            results = [t for t in results if t.column_id == column_id]
        return results

    def count_by_column(self, column_id: str) -> int:
        return sum(1 for t in self.tasks.values() if t.column_id == column_id)

    def update_column(self, task_id: str, new_column_id: str, updated_at: datetime) -> Optional[TaskRecord]:
        t = self.tasks.get(task_id)
        if not t:
            return None
        updated = TaskRecord(
            id=t.id,
            board_id=t.board_id,
            column_id=new_column_id,
            title=t.title,
            description=t.description,
            priority=t.priority,
            due_date=t.due_date,
            created_at=t.created_at,
            updated_at=updated_at,
        )
        self.tasks[task_id] = updated
        return updated

    def delete(self, task_id: str) -> bool:
        if task_id in self.tasks:
            del self.tasks[task_id]
            return True
        return False


def test_board_usecases():
    board_repo = InMemoryBoardRepository()
    create_uc = CreateBoardService(board_repo)
    get_uc = GetBoardService(board_repo)
    list_uc = ListBoardsService(board_repo)
    add_col_uc = AddColumnService(board_repo)

    # 1. Create board
    created = create_uc.exec(CreateBoardInput(name="Sprint 1", description="MVP release"))
    assert created.name == "Sprint 1"
    assert len(created.columns) == 3

    # 2. Get board
    retrieved = get_uc.exec(GetBoardInput(board_id=created.id))
    assert retrieved.id == created.id
    assert retrieved.name == "Sprint 1"

    # 3. List boards
    all_boards = list_uc.exec(ListBoardsInput())
    assert len(all_boards.boards) == 1

    # 4. Add column
    col_out = add_col_uc.exec(AddColumnInput(board_id=created.id, name="QA", wip_limit=2))
    assert col_out.column.name == "QA"
    assert col_out.column.wip_limit == 2

    # Verify duplicate column name rejected
    with pytest.raises(InboundDuplicateColumnNameError):
        add_col_uc.exec(AddColumnInput(board_id=created.id, name="QA"))


def test_task_usecases_and_wip_limit():
    board_repo = InMemoryBoardRepository()
    task_repo = InMemoryTaskRepository()

    create_board_uc = CreateBoardService(board_repo)
    add_col_uc = AddColumnService(board_repo)
    create_task_uc = CreateTaskService(board_repo, task_repo)
    move_task_uc = MoveTaskService(board_repo, task_repo)
    list_tasks_uc = ListTasksService(board_repo, task_repo)
    delete_task_uc = DeleteTaskService(task_repo)

    # Setup board with a column that has WIP limit of 2
    board = create_board_uc.exec(CreateBoardInput(name="Kanban Board"))
    qa_col = add_col_uc.exec(AddColumnInput(board_id=board.id, name="QA", wip_limit=2)).column

    # Create task 1 in Backlog
    t1 = create_task_uc.exec(
        CreateTaskInput(
            board_id=board.id,
            column="Backlog",
            title="Task 1",
            priority="HIGH",
        )
    )
    assert t1.column_name == "Backlog"

    # Create task 2 in Backlog
    t2 = create_task_uc.exec(
        CreateTaskInput(
            board_id=board.id,
            column="Backlog",
            title="Task 2",
            priority="MEDIUM",
        )
    )

    # Create task 3 in Backlog
    t3 = create_task_uc.exec(
        CreateTaskInput(
            board_id=board.id,
            column="Backlog",
            title="Task 3",
            priority="LOW",
        )
    )

    # Move task 1 to QA (WIP limit 2: currently 0 -> now 1)
    m1 = move_task_uc.exec(MoveTaskInput(task_id=t1.id, target_column_name="QA"))
    assert m1.column_name == "QA"

    # Move task 2 to QA (WIP limit 2: currently 1 -> now 2)
    m2 = move_task_uc.exec(MoveTaskInput(task_id=t2.id, target_column_name="QA"))
    assert m2.column_name == "QA"

    # Move task 3 to QA -> should be rejected with InboundMoveWipLimitExceededError!
    with pytest.raises(InboundMoveWipLimitExceededError) as exc_info:
        move_task_uc.exec(MoveTaskInput(task_id=t3.id, target_column_name="QA"))

    assert exc_info.value.code == "WIP_LIMIT_EXCEEDED"
    assert exc_info.value.current_count == 2
    assert exc_info.value.limit == 2

    # List tasks on board
    tasks_res = list_tasks_uc.exec(ListTasksInput(board_id=board.id))
    assert len(tasks_res.tasks) == 3

    # Delete task 1
    del_res = delete_task_uc.exec(DeleteTaskInput(task_id=t1.id))
    assert del_res.deleted is True

    # Now QA has 1 task; task 3 can be moved!
    m3 = move_task_uc.exec(MoveTaskInput(task_id=t3.id, target_column_name="QA"))
    assert m3.column_name == "QA"

