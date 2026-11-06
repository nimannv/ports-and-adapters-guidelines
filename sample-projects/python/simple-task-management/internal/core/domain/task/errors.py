from typing import Optional


class TaskDomainError(Exception):
    """Base domain exception for task bounded domain."""

    def __init__(self, message: str, code: str = "TASK_DOMAIN_ERROR"):
        super().__init__(message)
        self.message = message
        self.code = code


class TaskNotFoundError(TaskDomainError):
    def __init__(self, task_id: str):
        super().__init__(f"Task with id '{task_id}' was not found.", code="TASK_NOT_FOUND")
        self.task_id = task_id


class WipLimitExceededError(TaskDomainError):
    def __init__(
        self,
        column_name: str,
        column_id: str,
        current_count: int,
        limit: int,
    ):
        super().__init__(
            f"Cannot move task to '{column_name}': column limit of {limit} tasks has been reached.",
            code="WIP_LIMIT_EXCEEDED",
        )
        self.column_name = column_name
        self.column_id = column_id
        self.current_count = current_count
        self.limit = limit


class InvalidTaskMovementError(TaskDomainError):
    def __init__(self, message: str):
        super().__init__(message, code="INVALID_TASK_MOVEMENT")


class InvalidTaskError(TaskDomainError):
    def __init__(self, message: str):
        super().__init__(message, code="INVALID_TASK")

