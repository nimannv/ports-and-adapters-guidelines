from internal.core.application.ports.inbound.board.create_board import ColumnDTO
from internal.core.application.ports.inbound.board.get_board import (
    BoardNotFoundError,
    GetBoardInput,
    GetBoardOutput,
    GetBoardUseCase,
)
from internal.core.application.ports.outbound.board import BoardRepository


class GetBoardService(GetBoardUseCase):
    def __init__(self, board_repo: BoardRepository):
        self._board_repo = board_repo

    def exec(self, req: GetBoardInput) -> GetBoardOutput:
        record = self._board_repo.get_by_id(req.board_id)
        if not record:
            raise BoardNotFoundError(req.board_id)

        return GetBoardOutput(
            id=record.id,
            name=record.name,
            description=record.description,
            columns=[
                ColumnDTO(
                    id=c.id,
                    name=c.name,
                    order_index=c.order_index,
                    wip_limit=c.wip_limit,
                )
                for c in record.columns
            ],
            created_at=record.created_at,
            updated_at=record.updated_at,
        )

