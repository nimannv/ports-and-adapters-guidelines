from internal.core.application.ports.inbound.board.add_column import (
    AddColumnInput,
    AddColumnOutput,
    AddColumnUseCase,
    InboundBoardNotFoundError,
    InboundDuplicateColumnNameError,
    InboundInvalidWipLimitError,
)
from internal.core.application.ports.inbound.board.create_board import ColumnDTO
from internal.core.application.ports.outbound.board import BoardRepository, ColumnRecord
from internal.core.domain.board.entities import Board, Column
from internal.core.domain.board.errors import DuplicateColumnNameError, InvalidWipLimitError


class AddColumnService(AddColumnUseCase):
    def __init__(self, board_repo: BoardRepository):
        self._board_repo = board_repo

    def exec(self, req: AddColumnInput) -> AddColumnOutput:
        board_record = self._board_repo.get_by_id(req.board_id)
        if not board_record:
            raise InboundBoardNotFoundError(req.board_id)

        # Reconstruct domain board to enforce domain invariants
        domain_columns = [
            Column(
                id=c.id,
                name=c.name,
                order_index=c.order_index,
                wip_limit=c.wip_limit,
            )
            for c in board_record.columns
        ]
        domain_board = Board(
            id=board_record.id,
            name=board_record.name,
            description=board_record.description,
            columns=domain_columns,
            created_at=board_record.created_at,
            updated_at=board_record.updated_at,
        )

        try:
            new_col = domain_board.add_column(name=req.name, wip_limit=req.wip_limit)
        except DuplicateColumnNameError as e:
            raise InboundDuplicateColumnNameError(e.name, e.board_id) from e
        except InvalidWipLimitError as e:
            raise InboundInvalidWipLimitError(e.wip_limit) from e

        column_record = ColumnRecord(
            id=new_col.id,
            board_id=domain_board.id,
            name=new_col.name,
            order_index=new_col.order_index,
            wip_limit=new_col.wip_limit,
        )

        saved_column = self._board_repo.add_column(domain_board.id, column_record)

        return AddColumnOutput(
            board_id=domain_board.id,
            column=ColumnDTO(
                id=saved_column.id,
                name=saved_column.name,
                order_index=saved_column.order_index,
                wip_limit=saved_column.wip_limit,
            ),
        )

