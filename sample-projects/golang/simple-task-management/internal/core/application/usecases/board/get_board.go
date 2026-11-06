package board

import (
	"context"
	"errors"
	"fmt"

	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
)

var _ boardin.GetBoardUseCase = (*GetBoardService)(nil)

type GetBoardService struct {
	boardRepo boardout.Repository
}

func NewGetBoardService(boardRepo boardout.Repository) *GetBoardService {
	return &GetBoardService{boardRepo: boardRepo}
}

func (s *GetBoardService) Exec(ctx context.Context, in boardin.GetBoardInput) (boardin.GetBoardOutput, error) {
	rec, err := s.boardRepo.GetByID(ctx, in.BoardID)
	if err != nil {
		if errors.Is(err, boardout.ErrNotFound) {
			return boardin.GetBoardOutput{}, fmt.Errorf("%w: board %q", boardin.ErrBoardNotFound, in.BoardID)
		}
		return boardin.GetBoardOutput{}, fmt.Errorf("get board %q: %w", in.BoardID, err)
	}

	outCols := make([]boardin.ColumnDTO, len(rec.Columns))
	for i, c := range rec.Columns {
		outCols[i] = boardin.ColumnDTO{
			ID:         c.ID,
			Name:       c.Name,
			OrderIndex: c.OrderIndex,
			WipLimit:   c.WipLimit,
		}
	}

	return boardin.GetBoardOutput{
		ID:          rec.ID,
		Name:        rec.Name,
		Description: rec.Description,
		Columns:     outCols,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}, nil
}

