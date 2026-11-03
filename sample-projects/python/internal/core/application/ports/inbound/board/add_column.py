from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Optional

from internal.core.application.ports.inbound.board.create_board import ColumnDTO


@dataclass(frozen=True)
class AddColumnInput:
    board_id: str
    name: str
    wip_limit: Optional[int] = None


@dataclass(frozen=True)
class AddColumnOutput:
    board_id: str
    column: ColumnDTO


class InboundBoardNotFoundError(Exception):
    def __init__(self, board_id: str):
        super().__init__(f"Board with id '{board_id}' was not found.")
        self.code = "BOARD_NOT_FOUND"
        self.board_id = board_id


class InboundDuplicateColumnNameError(Exception):
    def __init__(self, name: str, board_id: str):
        super().__init__(f"A column with name '{name}' already exists on board '{board_id}'.")
        self.code = "DUPLICATE_COLUMN_NAME"
        self.name = name
        self.board_id = board_id


class InboundInvalidWipLimitError(Exception):
    def __init__(self, wip_limit: int):
        super().__init__(f"WIP limit must be greater than 0, got {wip_limit}.")
        self.code = "INVALID_WIP_LIMIT"
        self.wip_limit = wip_limit


class AddColumnUseCase(ABC):
    @abstractmethod
    def exec(self, req: AddColumnInput) -> AddColumnOutput:
        pass

