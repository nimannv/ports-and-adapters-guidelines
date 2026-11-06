package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
)

func writeUseCaseError(c *fiber.Ctx, err error) error {
	var wipErr *taskin.WipLimitExceededError
	switch {
	case errors.As(err, &wipErr):
		return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
			Error: ErrorDetail{
				Code:         "WIP_LIMIT_EXCEEDED",
				Message:      wipErr.Message,
				ColumnID:     wipErr.ColumnID,
				CurrentCount: wipErr.CurrentCount,
				Limit:        wipErr.Limit,
			},
		})
	case errors.Is(err, boardin.ErrBoardNotFound), errors.Is(err, taskin.ErrBoardNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "BOARD_NOT_FOUND", Message: err.Error()},
		})
	case errors.Is(err, boardin.ErrColumnNotFound), errors.Is(err, taskin.ErrColumnNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "COLUMN_NOT_FOUND", Message: err.Error()},
		})
	case errors.Is(err, taskin.ErrTaskNotFound):
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "TASK_NOT_FOUND", Message: err.Error()},
		})
	case errors.Is(err, boardin.ErrDuplicateColumnName):
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "DUPLICATE_COLUMN_NAME", Message: err.Error()},
		})
	case errors.Is(err, boardin.ErrInvalidWipLimit):
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_WIP_LIMIT", Message: err.Error()},
		})
	case errors.Is(err, boardin.ErrInvalidBoard):
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_BOARD", Message: err.Error()},
		})
	case errors.Is(err, taskin.ErrInvalidTaskMovement):
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_TASK_MOVEMENT", Message: err.Error()},
		})
	case errors.Is(err, taskin.ErrInvalidTask):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INVALID_TASK", Message: err.Error()},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "INTERNAL_ERROR", Message: "an internal server error occurred"},
		})
	}
}

