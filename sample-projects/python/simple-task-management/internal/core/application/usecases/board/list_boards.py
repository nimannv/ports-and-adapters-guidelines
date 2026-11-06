from internal.core.application.ports.inbound.board.create_board import ColumnDTO
from internal.core.application.ports.inbound.board.list_boards import (
    BoardSummaryDTO,
    ListBoardsInput,
    ListBoardsOutput,
    ListBoardsUseCase,
)
from internal.core.application.ports.outbound.board import BoardRepository


class ListBoardsService(ListBoardsUseCase):
    def __init__(self, board_repo: BoardRepository):
        self._board_repo = board_repo

    def exec(self, req: ListBoardsInput) -> ListBoardsOutput:
        records = self._board_repo.list_all()
        summaries = [
            BoardSummaryDTO(
                id=b.id,
                name=b.name,
                description=b.description,
                columns=[
                    ColumnDTO(
                        id=c.id,
                        name=c.name,
                        order_index=c.order_index,
                        wip_limit=c.wip_limit,
                    )
                    for c in b.columns
                ],
                created_at=b.created_at,
                updated_at=b.updated_at,
            )
            for b in records
        ]
        return ListBoardsOutput(boards=summaries)

