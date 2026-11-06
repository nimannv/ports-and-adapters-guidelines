package board

import (
	"context"
	"fmt"

	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
	boardout "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/outbound/board"
)

var _ boardin.ListBoardsUseCase = (*ListBoardsService)(nil)

type ListBoardsService struct {
	boardRepo boardout.Repository
}

func NewListBoardsService(boardRepo boardout.Repository) *ListBoardsService {
	return &ListBoardsService{boardRepo: boardRepo}
}

func (s *ListBoardsService) Exec(ctx context.Context, _ boardin.ListBoardsInput) (boardin.ListBoardsOutput, error) {
	records, err := s.boardRepo.ListAll(ctx)
	if err != nil {
		return boardin.ListBoardsOutput{}, fmt.Errorf("list boards: %w", err)
	}

	outBoards := make([]boardin.BoardItemDTO, len(records))
	for i, r := range records {
		cols := make([]boardin.ColumnDTO, len(r.Columns))
		for j, c := range r.Columns {
			cols[j] = boardin.ColumnDTO{
				ID:         c.ID,
				Name:       c.Name,
				OrderIndex: c.OrderIndex,
				WipLimit:   c.WipLimit,
			}
		}
		outBoards[i] = boardin.BoardItemDTO{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Columns:     cols,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}

	return boardin.ListBoardsOutput{Boards: outBoards}, nil
}

