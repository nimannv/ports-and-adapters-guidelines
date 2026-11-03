from internal.core.application.ports.inbound.board.create_board import (
    ColumnDTO,
    CreateBoardError,
    CreateBoardInput,
    CreateBoardOutput,
    CreateBoardUseCase,
)
from internal.core.application.ports.outbound.board import BoardRecord, BoardRepository, ColumnRecord
from internal.core.domain.board.entities import Board
from internal.core.domain.board.errors import DomainError


class CreateBoardService(CreateBoardUseCase):
    def __init__(self, board_repo: BoardRepository):
        self._board_repo = board_repo

    def exec(self, req: CreateBoardInput) -> CreateBoardOutput:
        try:
            domain_board = Board.create(
                name=req.name,
                description=req.description,
                custom_columns=req.custom_columns,
            )
        except DomainError as err:
            raise CreateBoardError(err.message, code=err.code) from err

        column_records = [
            ColumnRecord(
                id=c.id,
                board_id=domain_board.id,
                name=c.name,
                order_index=c.order_index,
                wip_limit=c.wip_limit,
            )
            for c in domain_board.columns
        ]

        board_record = BoardRecord(
            id=domain_board.id,
            name=domain_board.name,
            description=domain_board.description,
            columns=column_records,
            created_at=domain_board.created_at,
            updated_at=domain_board.updated_at,
        )

        saved = self._board_repo.save(board_record)

        return CreateBoardOutput(
            id=saved.id,
            name=saved.name,
            description=saved.description,
            columns=[
                ColumnDTO(
                    id=c.id,
                    name=c.name,
                    order_index=c.order_index,
                    wip_limit=c.wip_limit,
                )
                for c in saved.columns
            ],
            created_at=saved.created_at,
            updated_at=saved.updated_at,
        )

