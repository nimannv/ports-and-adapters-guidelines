package http

import "time"

// ColumnDTO represents a column in board responses.
type ColumnDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OrderIndex int    `json:"order_index"`
	WipLimit   *int   `json:"wip_limit,omitempty"`
}

// CustomColumnRequest is an optional column spec when creating a board.
type CustomColumnRequest struct {
	Name       string `json:"name"`
	OrderIndex int    `json:"order_index"`
	WipLimit   *int   `json:"wip_limit,omitempty"`
}

// CreateBoardRequest is the transport request body for POST /api/v1/boards.
type CreateBoardRequest struct {
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	CustomColumns []CustomColumnRequest `json:"custom_columns,omitempty"`
}

// BoardResponse is the transport response body for board objects.
type BoardResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Columns     []ColumnDTO `json:"columns"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// AddColumnRequest is the transport request body for POST /api/v1/boards/:boardId/columns.
type AddColumnRequest struct {
	Name     string `json:"name"`
	WipLimit *int   `json:"wip_limit,omitempty"`
}

// ColumnResponse is the response body for single column endpoints.
type ColumnResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	OrderIndex int    `json:"order_index"`
	WipLimit   *int   `json:"wip_limit,omitempty"`
}

// CreateTaskRequest is the transport request body for POST /api/v1/boards/:boardId/tasks.
type CreateTaskRequest struct {
	Column      string     `json:"column"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// TaskResponse is the transport response body for task items.
type TaskResponse struct {
	ID          string     `json:"id"`
	BoardID     string     `json:"board_id"`
	ColumnID    string     `json:"column_id"`
	ColumnName  string     `json:"column_name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// MoveTaskRequest is the transport request body for PATCH /api/v1/tasks/:taskId/move.
type MoveTaskRequest struct {
	TargetColumnID   string `json:"target_column_id,omitempty"`
	TargetColumnName string `json:"target_column_name,omitempty"`
}

// MoveTaskResponse is the transport response body for PATCH /api/v1/tasks/:taskId/move.
type MoveTaskResponse struct {
	TaskID     string    `json:"task_id"`
	BoardID    string    `json:"board_id"`
	ColumnID   string    `json:"column_id"`
	ColumnName string    `json:"column_name"`
	Title      string    `json:"title"`
	Priority   string    `json:"priority"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ErrorDetail describes an API error response.
type ErrorDetail struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	ColumnID     string `json:"column_id,omitempty"`
	CurrentCount int    `json:"current_count,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

// ErrorResponse wraps the standard API error structure.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

