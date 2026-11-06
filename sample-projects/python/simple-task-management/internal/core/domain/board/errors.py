from typing import Optional


class DomainError(Exception):
    """Base domain exception."""

    def __init__(self, message: str, code: str = "DOMAIN_ERROR"):
        super().__init__(message)
        self.message = message
        self.code = code


class BoardNotFoundError(DomainError):
    def __init__(self, board_id: str):
        super().__init__(f"Board with id '{board_id}' was not found.", code="BOARD_NOT_FOUND")
        self.board_id = board_id


class ColumnNotFoundError(DomainError):
    def __init__(self, column_ref: str, board_id: Optional[str] = None):
        msg = f"Column '{column_ref}' was not found"
        if board_id:
            msg += f" on board '{board_id}'."
        else:
            msg += "."
        super().__init__(msg, code="COLUMN_NOT_FOUND")
        self.column_ref = column_ref
        self.board_id = board_id


class DuplicateColumnNameError(DomainError):
    def __init__(self, name: str, board_id: str):
        super().__init__(
            f"A column with name '{name}' already exists on board '{board_id}'.",
            code="DUPLICATE_COLUMN_NAME",
        )
        self.name = name
        self.board_id = board_id


class ColumnNotEmptyError(DomainError):
    def __init__(self, column_id: str, task_count: int):
        super().__init__(
            f"Cannot delete column '{column_id}' because it still contains {task_count} tasks.",
            code="COLUMN_NOT_EMPTY",
        )
        self.column_id = column_id
        self.task_count = task_count


class InvalidBoardError(DomainError):
    def __init__(self, message: str):
        super().__init__(message, code="INVALID_BOARD")


class InvalidWipLimitError(DomainError):
    def __init__(self, wip_limit: int):
        super().__init__(
            f"WIP limit must be greater than 0, got {wip_limit}.",
            code="INVALID_WIP_LIMIT",
        )
        self.wip_limit = wip_limit

