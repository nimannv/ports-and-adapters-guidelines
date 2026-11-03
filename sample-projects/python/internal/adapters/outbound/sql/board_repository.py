from typing import List, Optional
from sqlalchemy.orm import sessionmaker

from internal.adapters.outbound.sql.mappers import board_model_to_record, column_model_to_record
from internal.adapters.outbound.sql.models import BoardModel, ColumnModel
from internal.core.application.ports.outbound.board import (
    BoardRecord,
    BoardRepository,
    ColumnRecord,
    OutboundRepositoryError,
)


class SqlAlchemyBoardRepository(BoardRepository):
    def __init__(self, session_factory: sessionmaker):
        self._session_factory = session_factory

    def save(self, record: BoardRecord) -> BoardRecord:
        try:
            with self._session_factory() as session:
                existing = session.query(BoardModel).filter(BoardModel.id == record.id).first()
                if not existing:
                    board_model = BoardModel(
                        id=record.id,
                        name=record.name,
                        description=record.description,
                        created_at=record.created_at,
                        updated_at=record.updated_at,
                    )
                    for c in record.columns:
                        col_model = ColumnModel(
                            id=c.id,
                            board_id=record.id,
                            name=c.name,
                            order_index=c.order_index,
                            wip_limit=c.wip_limit,
                        )
                        board_model.columns.append(col_model)
                    session.add(board_model)
                else:
                    existing.name = record.name
                    existing.description = record.description
                    existing.updated_at = record.updated_at
                session.commit()

                # Refresh and retrieve
                saved = session.query(BoardModel).filter(BoardModel.id == record.id).first()
                return board_model_to_record(saved)
        except Exception as e:
            raise OutboundRepositoryError(f"Database error while saving board: {e}") from e

    def get_by_id(self, board_id: str) -> Optional[BoardRecord]:
        try:
            with self._session_factory() as session:
                model = session.query(BoardModel).filter(BoardModel.id == board_id).first()
                if not model:
                    return None
                return board_model_to_record(model)
        except Exception as e:
            raise OutboundRepositoryError(f"Database error retrieving board '{board_id}': {e}") from e

    def list_all(self) -> List[BoardRecord]:
        try:
            with self._session_factory() as session:
                models = session.query(BoardModel).order_by(BoardModel.created_at.desc()).all()
                return [board_model_to_record(m) for m in models]
        except Exception as e:
            raise OutboundRepositoryError(f"Database error listing boards: {e}") from e

    def add_column(self, board_id: str, column: ColumnRecord) -> ColumnRecord:
        try:
            with self._session_factory() as session:
                col_model = ColumnModel(
                    id=column.id,
                    board_id=board_id,
                    name=column.name,
                    order_index=column.order_index,
                    wip_limit=column.wip_limit,
                )
                session.add(col_model)
                session.commit()
                return column_model_to_record(col_model)
        except Exception as e:
            raise OutboundRepositoryError(f"Database error adding column: {e}") from e

