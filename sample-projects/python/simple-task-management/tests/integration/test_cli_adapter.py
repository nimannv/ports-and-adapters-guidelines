from typer.testing import CliRunner

from internal.adapters.inbound.cli import create_cli_app
from internal.adapters.outbound.sql import (
    SqlAlchemyBoardRepository,
    SqlAlchemyTaskRepository,
    create_db_engine,
    create_session_factory,
    init_db,
)
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


def create_test_cli():
    engine = create_db_engine("sqlite:///:memory:")
    init_db(engine)
    session_factory = create_session_factory(engine)

    board_repo = SqlAlchemyBoardRepository(session_factory)
    task_repo = SqlAlchemyTaskRepository(session_factory)

    return create_cli_app(
        create_board_uc=CreateBoardService(board_repo),
        get_board_uc=GetBoardService(board_repo),
        list_boards_uc=ListBoardsService(board_repo),
        add_column_uc=AddColumnService(board_repo),
        create_task_uc=CreateTaskService(board_repo, task_repo),
        move_task_uc=MoveTaskService(board_repo, task_repo),
        list_tasks_uc=ListTasksService(board_repo, task_repo),
        delete_task_uc=DeleteTaskService(task_repo),
    )


def test_cli_lifecycle():
    app = create_test_cli()
    runner = CliRunner()

    # 1. Create board
    res = runner.invoke(app, ["board", "create", "--name", "Sprint 1", "--description", "Release MVP"])
    assert res.exit_code == 0
    assert "Board created:" in res.stdout
    assert "Sprint 1" in res.stdout
    assert "Backlog (unlimited)" in res.stdout

    # Extract Board ID
    lines = res.stdout.splitlines()
    board_id_line = [l for l in lines if "ID:" in l][0]
    board_id = board_id_line.split("ID:")[1].strip()

    # 2. Add column with WIP limit
    res = runner.invoke(app, ["column", "add", "--board-id", board_id, "--name", "QA Review", "--wip-limit", "2"])
    assert res.exit_code == 0
    assert "Column 'QA Review' added with WIP limit 2." in res.stdout

    # 3. Create task
    res = runner.invoke(
        app,
        [
            "task",
            "create",
            "--board-id",
            board_id,
            "--column",
            "Backlog",
            "--title",
            "Implement Login",
            "--priority",
            "HIGH",
            "--due-date",
            "2026-12-01",
        ],
    )
    assert res.exit_code == 0
    assert "Task created: ID" in res.stdout
    assert "[HIGH]" in res.stdout

    # Extract Task ID
    task_id = res.stdout.split("Task created: ID ")[1].split(" ")[0].strip()

    # 4. Move task
    res = runner.invoke(
        app,
        [
            "task",
            "move",
            "--task-id",
            task_id,
            "--to-column",
            "In Progress",
        ],
    )
    assert res.exit_code == 0
    assert f"Task {task_id} moved to 'In Progress'." in res.stdout

    # 5. List tasks
    res = runner.invoke(app, ["task", "list", "--board-id", board_id])
    assert res.exit_code == 0
    assert "Implement Login" in res.stdout
    assert "In Progress" in res.stdout
    assert "HIGH" in res.stdout

    # 6. Delete task
    res = runner.invoke(app, ["task", "delete", "--task-id", task_id])
    assert res.exit_code == 0
    assert f"Task {task_id} deleted successfully." in res.stdout

