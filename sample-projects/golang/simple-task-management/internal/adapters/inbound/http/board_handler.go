package http

import (
	"github.com/gofiber/fiber/v2"
	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
)

type BoardHandler struct {
	createBoard boardin.CreateBoardUseCase
	getBoard    boardin.GetBoardUseCase
	listBoards  boardin.ListBoardsUseCase
	addColumn   boardin.AddColumnUseCase
	createTask  taskin.CreateTaskUseCase
	listTasks   taskin.ListTasksUseCase
}

func NewBoardHandler(
	createBoard boardin.CreateBoardUseCase,
	getBoard boardin.GetBoardUseCase,
	listBoards boardin.ListBoardsUseCase,
	addColumn boardin.AddColumnUseCase,
	createTask taskin.CreateTaskUseCase,
	listTasks taskin.ListTasksUseCase,
) *BoardHandler {
	return &BoardHandler{
		createBoard: createBoard,
		getBoard:    getBoard,
		listBoards:  listBoards,
		addColumn:   addColumn,
		createTask:  createTask,
		listTasks:   listTasks,
	}
}

func (h *BoardHandler) Register(router fiber.Router) {
	boards := router.Group("/boards")

	boards.Post("/", h.handleCreateBoard)
	boards.Get("/", h.handleListBoards)
	boards.Get("/:boardId", h.handleGetBoard)
	boards.Post("/:boardId/columns", h.handleAddColumn)
	boards.Post("/:boardId/tasks", h.handleCreateTask)
	boards.Get("/:boardId/tasks", h.handleListTasks)
}

func (h *BoardHandler) handleCreateBoard(c *fiber.Ctx) error {
	var req CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "MALFORMED_JSON", Message: "cannot parse request body"},
		})
	}

	var customCols []boardin.CustomColumnInput
	for _, col := range req.CustomColumns {
		customCols = append(customCols, boardin.CustomColumnInput{
			Name:       col.Name,
			OrderIndex: col.OrderIndex,
			WipLimit:   col.WipLimit,
		})
	}

	out, err := h.createBoard.Exec(c.UserContext(), boardin.CreateBoardInput{
		Name:          req.Name,
		Description:   req.Description,
		CustomColumns: customCols,
	})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	cols := make([]ColumnDTO, len(out.Columns))
	for i, col := range out.Columns {
		cols[i] = ColumnDTO{
			ID:         col.ID,
			Name:       col.Name,
			OrderIndex: col.OrderIndex,
			WipLimit:   col.WipLimit,
		}
	}

	return c.Status(fiber.StatusCreated).JSON(BoardResponse{
		ID:          out.ID,
		Name:        out.Name,
		Description: out.Description,
		Columns:     cols,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	})
}

func (h *BoardHandler) handleListBoards(c *fiber.Ctx) error {
	out, err := h.listBoards.Exec(c.UserContext(), boardin.ListBoardsInput{})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	boards := make([]BoardResponse, len(out.Boards))
	for i, b := range out.Boards {
		cols := make([]ColumnDTO, len(b.Columns))
		for j, col := range b.Columns {
			cols[j] = ColumnDTO{
				ID:         col.ID,
				Name:       col.Name,
				OrderIndex: col.OrderIndex,
				WipLimit:   col.WipLimit,
			}
		}
		boards[i] = BoardResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			Columns:     cols,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(boards)
}

func (h *BoardHandler) handleGetBoard(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	out, err := h.getBoard.Exec(c.UserContext(), boardin.GetBoardInput{BoardID: boardID})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	cols := make([]ColumnDTO, len(out.Columns))
	for i, col := range out.Columns {
		cols[i] = ColumnDTO{
			ID:         col.ID,
			Name:       col.Name,
			OrderIndex: col.OrderIndex,
			WipLimit:   col.WipLimit,
		}
	}

	return c.Status(fiber.StatusOK).JSON(BoardResponse{
		ID:          out.ID,
		Name:        out.Name,
		Description: out.Description,
		Columns:     cols,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	})
}

func (h *BoardHandler) handleAddColumn(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	var req AddColumnRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "MALFORMED_JSON", Message: "cannot parse request body"},
		})
	}

	out, err := h.addColumn.Exec(c.UserContext(), boardin.AddColumnInput{
		BoardID:  boardID,
		Name:     req.Name,
		WipLimit: req.WipLimit,
	})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(ColumnResponse{
		ID:         out.Column.ID,
		Name:       out.Column.Name,
		OrderIndex: out.Column.OrderIndex,
		WipLimit:   out.Column.WipLimit,
	})
}

func (h *BoardHandler) handleCreateTask(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	var req CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: ErrorDetail{Code: "MALFORMED_JSON", Message: "cannot parse request body"},
		})
	}

	out, err := h.createTask.Exec(c.UserContext(), taskin.CreateTaskInput{
		BoardID:     boardID,
		Column:      req.Column,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(TaskResponse{
		ID:          out.ID,
		BoardID:     out.BoardID,
		ColumnID:    out.ColumnID,
		ColumnName:  out.ColumnName,
		Title:       out.Title,
		Description: out.Description,
		Priority:    out.Priority,
		DueDate:     out.DueDate,
		CreatedAt:   out.CreatedAt,
		UpdatedAt:   out.UpdatedAt,
	})
}

func (h *BoardHandler) handleListTasks(c *fiber.Ctx) error {
	boardID := c.Params("boardId")
	columnID := c.Query("column_id")

	out, err := h.listTasks.Exec(c.UserContext(), taskin.ListTasksInput{
		BoardID:  boardID,
		ColumnID: columnID,
	})
	if err != nil {
		return writeUseCaseError(c, err)
	}

	tasks := make([]TaskResponse, len(out.Tasks))
	for i, t := range out.Tasks {
		tasks[i] = TaskResponse{
			ID:          t.ID,
			BoardID:     t.BoardID,
			ColumnID:    t.ColumnID,
			ColumnName:  t.ColumnName,
			Title:       t.Title,
			Description: t.Description,
			Priority:    t.Priority,
			DueDate:     t.DueDate,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}

