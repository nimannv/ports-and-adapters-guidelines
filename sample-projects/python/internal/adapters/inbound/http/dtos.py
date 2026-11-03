from datetime import datetime
from typing import Any, Dict, List, Optional
from pydantic import BaseModel, Field


class CreateColumnPayload(BaseModel):
    name: str
    order_index: Optional[int] = None
    wip_limit: Optional[int] = None


class CreateBoardRequest(BaseModel):
    name: str = Field(..., min_length=1)
    description: str = ""
    custom_columns: Optional[List[CreateColumnPayload]] = None


class AddColumnRequest(BaseModel):
    name: str = Field(..., min_length=1)
    wip_limit: Optional[int] = None


class CreateTaskRequest(BaseModel):
    column: str = Field(..., min_length=1, description="Column ID or Name")
    title: str = Field(..., min_length=1)
    description: str = ""
    priority: str = "MEDIUM"
    due_date: Optional[datetime] = None


class MoveTaskRequest(BaseModel):
    target_column_id: Optional[str] = None
    target_column_name: Optional[str] = None


class ColumnResponse(BaseModel):
    id: str
    name: str
    order_index: int
    wip_limit: Optional[int] = None


class BoardResponse(BaseModel):
    id: str
    name: str
    description: str
    columns: List[ColumnResponse]
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class TaskResponse(BaseModel):
    id: str
    board_id: str
    column_id: str
    column_name: str
    title: str
    description: str
    priority: str
    due_date: Optional[datetime] = None
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class MoveTaskResponse(BaseModel):
    task_id: str
    board_id: str
    column_id: str
    column_name: str
    title: str
    priority: str
    updated_at: Optional[datetime] = None


class ErrorDetail(BaseModel):
    code: str
    message: str
    column_id: Optional[str] = None
    current_count: Optional[int] = None
    limit: Optional[int] = None


class ErrorResponse(BaseModel):
    error: Dict[str, Any]

