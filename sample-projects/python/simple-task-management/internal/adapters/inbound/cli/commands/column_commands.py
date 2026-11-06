from typing import Optional
import typer
from rich.console import Console

from internal.core.application.ports.inbound.board import (
    AddColumnInput,
    AddColumnUseCase,
    BoardNotFoundError,
    InboundDuplicateColumnNameError,
    InboundInvalidWipLimitError,
)

console = Console()


def create_column_subapp(add_column_uc: AddColumnUseCase) -> typer.Typer:
    app = typer.Typer(help="Manage board workflow columns")

    @app.command("add")
    def add(
        board_id: str = typer.Option(..., "--board-id", "-b", help="UUID of the board"),
        name: str = typer.Option(..., "--name", "-n", help="Name of the new column"),
        wip_limit: Optional[int] = typer.Option(None, "--wip-limit", "-w", help="Optional WIP limit (>0)"),
    ):
        """Add a column with an optional WIP limit to an existing board."""
        try:
            res = add_column_uc.exec(
                AddColumnInput(board_id=board_id, name=name, wip_limit=wip_limit)
            )
            col = res.column
            if col.wip_limit is not None:
                console.print(f"Column '{col.name}' added with WIP limit {col.wip_limit}.")
            else:
                console.print(f"Column '{col.name}' added with unlimited capacity.")
        except BoardNotFoundError:
            console.print(f"[bold red]Error:[/bold red] Board '{board_id}' not found.")
            raise typer.Exit(code=1)
        except InboundDuplicateColumnNameError as e:
            console.print(f"[bold red]Error:[/bold red] Column '{e.name}' already exists on board '{e.board_id}'.")
            raise typer.Exit(code=1)
        except InboundInvalidWipLimitError as e:
            console.print(f"[bold red]Error:[/bold red] WIP limit must be greater than 0, got {e.wip_limit}.")
            raise typer.Exit(code=1)
        except Exception as e:
            console.print(f"[bold red]Error:[/bold red] {e}")
            raise typer.Exit(code=1)

    return app

