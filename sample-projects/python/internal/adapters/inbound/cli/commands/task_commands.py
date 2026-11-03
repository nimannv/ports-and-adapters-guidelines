from datetime import datetime
from typing import Optional
import typer
from rich.console import Console

from internal.adapters.inbound.cli.formatters import print_tasks_table
from internal.core.application.ports.inbound.task import (
    CreateTaskInput,
    CreateTaskUseCase,
    DeleteTaskInput,
    DeleteTaskUseCase,
    InboundBoardNotFoundError,
    InboundColumnNotFoundError,
    InboundDeleteTaskNotFoundError,
    InboundInvalidTaskError,
    InboundInvalidTaskMovementError,
    InboundListTasksBoardNotFoundError,
    InboundMoveColumnNotFoundError,
    InboundMoveWipLimitExceededError,
    InboundTaskNotFoundError,
    InboundWipLimitExceededError,
    ListTasksInput,
    ListTasksUseCase,
    MoveTaskInput,
    MoveTaskUseCase,
)

console = Console()


def create_task_subapp(
    create_task_uc: CreateTaskUseCase,
    move_task_uc: MoveTaskUseCase,
    list_tasks_uc: ListTasksUseCase,
    delete_task_uc: DeleteTaskUseCase,
) -> typer.Typer:
    app = typer.Typer(help="Manage tasks")

    @app.command("create")
    def create(
        board_id: str = typer.Option(..., "--board-id", "-b", help="UUID of the board"),
        column: str = typer.Option(..., "--column", "-c", help="Target column name or UUID"),
        title: str = typer.Option(..., "--title", "-t", help="Title of the task"),
        description: str = typer.Option("", "--description", "-d", help="Description of the task"),
        priority: str = typer.Option("MEDIUM", "--priority", "-p", help="LOW, MEDIUM, HIGH, URGENT"),
        due_date: Optional[str] = typer.Option(None, "--due-date", help="Due date (YYYY-MM-DD)"),
    ):
        """Create a new task in a specified column."""
        parsed_due_date = None
        if due_date:
            try:
                parsed_due_date = datetime.fromisoformat(due_date)
            except ValueError:
                console.print(f"[bold red]Error:[/bold red] Invalid due-date format '{due_date}'. Use YYYY-MM-DD.")
                raise typer.Exit(code=1)

        try:
            res = create_task_uc.exec(
                CreateTaskInput(
                    board_id=board_id,
                    column=column,
                    title=title,
                    description=description,
                    priority=priority.upper(),
                    due_date=parsed_due_date,
                )
            )
            console.print(f"Task created: ID {res.id} [{res.priority}]")
        except InboundBoardNotFoundError:
            console.print(f"[bold red]Error:[/bold red] Board '{board_id}' not found.")
            raise typer.Exit(code=1)
        except InboundColumnNotFoundError as e:
            console.print(f"[bold red]Error:[/bold red] Column '{e.column_ref}' not found on board '{board_id}'.")
            raise typer.Exit(code=1)
        except InboundWipLimitExceededError as e:
            console.print(f"[bold red]Error (WIP_LIMIT_EXCEEDED):[/bold red] {e.message}")
            raise typer.Exit(code=1)
        except InboundInvalidTaskError as e:
            console.print(f"[bold red]Error:[/bold red] {e.message}")
            raise typer.Exit(code=1)

    @app.command("move")
    def move(
        task_id: str = typer.Option(..., "--task-id", "-t", help="UUID of the task"),
        to_column: Optional[str] = typer.Option(None, "--to-column", "-c", help="Target column name"),
        to_column_id: Optional[str] = typer.Option(None, "--to-column-id", help="Target column UUID"),
    ):
        """Move a task to a different column on the same board."""
        if not to_column and not to_column_id:
            console.print("[bold red]Error:[/bold red] Either --to-column or --to-column-id must be provided.")
            raise typer.Exit(code=1)

        try:
            res = move_task_uc.exec(
                MoveTaskInput(
                    task_id=task_id,
                    target_column_id=to_column_id,
                    target_column_name=to_column,
                )
            )
            console.print(f"Task {res.task_id} moved to '{res.column_name}'.")
        except InboundTaskNotFoundError:
            console.print(f"[bold red]Error:[/bold red] Task '{task_id}' not found.")
            raise typer.Exit(code=1)
        except InboundMoveColumnNotFoundError as e:
            console.print(f"[bold red]Error:[/bold red] Target column '{e.column_ref}' not found.")
            raise typer.Exit(code=1)
        except InboundMoveWipLimitExceededError as e:
            console.print(f"[bold red]Error (WIP_LIMIT_EXCEEDED):[/bold red] {e.message}")
            raise typer.Exit(code=1)
        except InboundInvalidTaskMovementError as e:
            console.print(f"[bold red]Error:[/bold red] {e.message}")
            raise typer.Exit(code=1)

    @app.command("list")
    def list_tasks(
        board_id: str = typer.Option(..., "--board-id", "-b", help="UUID of the board"),
        column: Optional[str] = typer.Option(None, "--column", "-c", help="Filter by column name"),
        column_id: Optional[str] = typer.Option(None, "--column-id", help="Filter by column UUID"),
    ):
        """List tasks on a board, with optional column filter."""
        try:
            res = list_tasks_uc.exec(
                ListTasksInput(
                    board_id=board_id,
                    column_id=column_id,
                    column_name=column,
                )
            )
            if not res.tasks:
                console.print("[dim]No tasks found on this board.[/dim]")
                return
            print_tasks_table(res.tasks)
        except InboundListTasksBoardNotFoundError:
            console.print(f"[bold red]Error:[/bold red] Board '{board_id}' not found.")
            raise typer.Exit(code=1)

    @app.command("delete")
    def delete(
        task_id: str = typer.Option(..., "--task-id", "-t", help="UUID of the task to delete"),
    ):
        """Delete a task by ID."""
        try:
            delete_task_uc.exec(DeleteTaskInput(task_id=task_id))
            console.print(f"Task {task_id} deleted successfully.")
        except InboundDeleteTaskNotFoundError:
            console.print(f"[bold red]Error:[/bold red] Task '{task_id}' not found.")
            raise typer.Exit(code=1)

    return app

