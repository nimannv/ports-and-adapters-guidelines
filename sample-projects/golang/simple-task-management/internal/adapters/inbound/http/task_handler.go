package http

import (
	"github.com/gofiber/fiber/v2"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
)

type TaskHandler struct {
	moveTask   taskin.MoveTaskUseCase
	deleteTask taskin.DeleteTaskUseCase
}

func NewTaskHandler(
	moveTask taskin.MoveTaskUseCase,
	deleteTask taskin.DeleteTaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		moveTask:   moveTask,
		deleteTask: deleteTask,
	}
}

func (h *TaskHandler) Register(router fiber.Router) {
	tasks := router.Group("/tasks")

	tasks.Patch("/:taskId/move", h.handleMoveTask)
	tasks.Delete("/:taskId", h.handleDeleteTask)
}

func (h *TaskHandler) handleMoveTask(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	var req MoveTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "MALFORMED_JSON", Message: "cannot parse request body"},
		})
	}

	out, err := h.moveTask.Exec(c.UserContext(), taskin.MoveTaskInput{
		TaskID:           taskID,
		TargetColumnID:   req.TargetColumnID,
		TargetColumnName: req.TargetColumnName,
	})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(MoveTaskResponse{
		TaskID:     out.TaskID,
		BoardID:    out.BoardID,
		ColumnID:   out.ColumnID,
		ColumnName: out.ColumnName,
		Title:      out.Title,
		Priority:   out.Priority,
		UpdatedAt:  out.UpdatedAt,
	})
}

func (h *TaskHandler) handleDeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	_, err := h.deleteTask.Exec(c.UserContext(), taskin.DeleteTaskInput{
		TaskID: taskID,
	})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

