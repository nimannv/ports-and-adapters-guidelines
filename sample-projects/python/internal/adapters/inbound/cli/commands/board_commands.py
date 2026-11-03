from typing import Optional
import typer
from rich.console import Console
from rich.table import Table

from internal.adapters.inbound.cli.formatters import format_column_summary, print_board_created
from internal.core.application.ports.inbound.board import (
    BoardNotFoundError,
    CreateBoardError,
    CreateBoardInput,
    CreateBoardUseCase,
    GetBoardInput,
    GetBoardUseCase,
    ListBoardsInput,
    ListBoardsUseCase,
)

console = Console()


def create_board_subapp(
    create_board_uc: CreateBoardUseCase,
    get_board_uc: GetBoardUseCase,
    list_boards_uc: ListBoardsUseCase,
) -> typer.Typer:
    app = typer.Typer(help="Manage Kanban boards")

    @app.command("create")
    def create(
        name: str = typer.Option(..., "--name", "-n", help="Name of the board"),
        description: str = typer.Option("", "--description", "-d", help="Description of the board"),
    ):
        """Create a new board with default columns (Backlog, In Progress [limit: 3], Done)."""
        try:
            res = create_board_uc.exec(CreateBoardInput(name=name, description=description))
            print_board_created(res.id, res.name, res.columns)
        except CreateBoardError as e:
            console.print(f"[bold red]Error:[/bold red] {e.message}")
            raise typer.Exit(code=1)

    @app.command("list")
    def list_all():
        """List all boards."""
        res = list_boards_uc.exec(ListBoardsInput())
        if not res.boards:
            console.print("[dim]No boards found.[/dim]")
            return

        table = Table(title="Boards", show_header=True, header_style="bold blue")
        table.add_column("ID", style="cyan")
        table.add_column("NAME", style="white")
        table.add_column("COLUMNS", style="yellow")

        for b in res.boards:
            table.add_row(b.id, b.name, format_column_summary(b.columns))

        console.print(table)

    @app.command("get")
    def get(
        board_id: str = typer.Option(..., "--board-id", "-b", help="UUID of the board"),
    ):
        """Get board details."""
        try:
            b = get_board_uc.exec(GetBoardInput(board_id=board_id))
            console.print(f"[bold]Board:[/bold] {b.name} ({b.id})")
            if b.description:
                console.print(f"[dim]{b.description}[/dim]")
            console.print(f"Columns: {format_column_summary(b.columns)}")
        except BoardNotFoundError as e:
            console.print(f"[bold red]Error:[/bold red] Board '{board_id}' not found.")
            raise typer.Exit(code=1)

    return app

