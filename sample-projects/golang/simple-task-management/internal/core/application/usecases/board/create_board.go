package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
	boarddomain "github.com/ports-and-adapters/simple-task-management/internal/core/domain/board"
)

var _ boardin.CreateBoardUseCase = (*CreateBoardService)(nil)

type CreateBoardService struct {
	boardRepo boardout.Repository
}

func NewCreateBoardService(boardRepo boardout.Repository) *CreateBoardService {
	return &CreateBoardService{boardRepo: boardRepo}
}

func (s *CreateBoardService) Exec(ctx context.Context, in boardin.CreateBoardInput) (boardin.CreateBoardOutput, error) {
	var domainCols []boarddomain.Column
	if len(in.CustomColumns) > 0 {
		for i, c := range in.CustomColumns {
			order := c.OrderIndex
			if order == 0 && i > 0 {
				order = i
			}
			col, err := boarddomain.NewColumn("", c.Name, order, c.WipLimit)
			if err != nil {
				return boardin.CreateBoardOutput{}, translateDomainError(err)
			}
			domainCols = append(domainCols, col)
		}
	}

	b, err := boarddomain.NewBoard("", in.Name, in.Description, domainCols, time.Now().UTC())
	if err != nil {
		return boardin.CreateBoardOutput{}, translateDomainError(err)
	}

	// Map domain entity to outbound DTO
	recCols := make([]boardout.ColumnRecord, len(b.Columns()))
	for i, c := range b.Columns() {
		recCols[i] = boardout.ColumnRecord{
			ID:         c.ID(),
			BoardID:    b.ID(),
			Name:       c.Name(),
			OrderIndex: c.OrderIndex(),
			WipLimit:   c.WipLimit(),
		}
	}

	rec := boardout.BoardRecord{
		ID:          b.ID(),
		Name:        b.Name(),
		Description: b.Description(),
		Columns:     recCols,
		CreatedAt:   b.CreatedAt(),
		UpdatedAt:   b.UpdatedAt(),
	}

	saved, err := s.boardRepo.Save(ctx, rec)
	if err != nil {
		if errors.Is(err, boardout.ErrDuplicate) {
			return boardin.CreateBoardOutput{}, fmt.Errorf("%w: %v", boardin.ErrDuplicateColumnName, err)
		}
		return boardin.CreateBoardOutput{}, fmt.Errorf("save board: %w", err)
	}

	outCols := make([]boardin.ColumnDTO, len(saved.Columns))
	for i, c := range saved.Columns {
		outCols[i] = boardin.ColumnDTO{
			ID:         c.ID,
			Name:       c.Name,
			OrderIndex: c.OrderIndex,
			WipLimit:   c.WipLimit,
		}
	}

	return boardin.CreateBoardOutput{
		ID:          saved.ID,
		Name:        saved.Name,
		Description: saved.Description,
		Columns:     outCols,
		CreatedAt:   saved.CreatedAt,
		UpdatedAt:   saved.UpdatedAt,
	}, nil
}

func translateDomainError(err error) error {
	switch {
	case errors.Is(err, boarddomain.ErrInvalidBoard):
		return fmt.Errorf("%w: %v", boardin.ErrInvalidBoard, err)
	case errors.Is(err, boarddomain.ErrDuplicateColumnName):
		return fmt.Errorf("%w: %v", boardin.ErrDuplicateColumnName, err)
	case errors.Is(err, boarddomain.ErrInvalidWipLimit):
		return fmt.Errorf("%w: %v", boardin.ErrInvalidWipLimit, err)
	case errors.Is(err, boarddomain.ErrColumnNotFound):
		return fmt.Errorf("%w: %v", boardin.ErrColumnNotFound, err)
	default:
		return err
	}
}

