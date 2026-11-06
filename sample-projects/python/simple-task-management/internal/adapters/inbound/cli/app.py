import typer

from internal.adapters.inbound.cli.commands.board_commands import create_board_subapp
from internal.adapters.inbound.cli.commands.column_commands import create_column_subapp
from internal.adapters.inbound.cli.commands.task_commands import create_task_subapp
from internal.core.application.ports.inbound.board import (
    AddColumnUseCase,
    CreateBoardUseCase,
    GetBoardUseCase,
    ListBoardsUseCase,
)
from internal.core.application.ports.inbound.task import (
    CreateTaskUseCase,
    DeleteTaskUseCase,
    ListTasksUseCase,
    MoveTaskUseCase,
)


def create_cli_app(
    create_board_uc: CreateBoardUseCase,
    get_board_uc: GetBoardUseCase,
    list_boards_uc: ListBoardsUseCase,
    add_column_uc: AddColumnUseCase,
    create_task_uc: CreateTaskUseCase,
    move_task_uc: MoveTaskUseCase,
    list_tasks_uc: ListTasksUseCase,
    delete_task_uc: DeleteTaskUseCase,
) -> typer.Typer:
    app = typer.Typer(
        name="task-cli",
        help="Kanban-style Simple Task Management System CLI",
        no_args_is_help=True,
    )

    board_subapp = create_board_subapp(create_board_uc, get_board_uc, list_boards_uc)
    column_subapp = create_column_subapp(add_column_uc)
    task_subapp = create_task_subapp(create_task_uc, move_task_uc, list_tasks_uc, delete_task_uc)

    app.add_typer(board_subapp, name="board")
    app.add_typer(column_subapp, name="column")
    app.add_typer(task_subapp, name="task")

    return app

