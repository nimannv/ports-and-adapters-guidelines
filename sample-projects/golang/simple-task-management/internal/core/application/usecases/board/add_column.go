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

var _ boardin.AddColumnUseCase = (*AddColumnService)(nil)

type AddColumnService struct {
	boardRepo boardout.Repository
}

func NewAddColumnService(boardRepo boardout.Repository) *AddColumnService {
	return &AddColumnService{boardRepo: boardRepo}
}

func (s *AddColumnService) Exec(ctx context.Context, in boardin.AddColumnInput) (boardin.AddColumnOutput, error) {
	rec, err := s.boardRepo.GetByID(ctx, in.BoardID)
	if err != nil {
		if errors.Is(err, boardout.ErrNotFound) {
			return boardin.AddColumnOutput{}, fmt.Errorf("%w: board %q", boardin.ErrBoardNotFound, in.BoardID)
		}
		return boardin.AddColumnOutput{}, fmt.Errorf("get board %q: %w", in.BoardID, err)
	}

	// Rehydrate domain board to check invariants
	domainCols := make([]boarddomain.Column, len(rec.Columns))
	for i, c := range rec.Columns {
		domainCols[i] = boarddomain.RehydrateColumn(c.ID, c.Name, c.OrderIndex, c.WipLimit)
	}
	b := boarddomain.RehydrateBoard(rec.ID, rec.Name, rec.Description, domainCols, rec.CreatedAt, rec.UpdatedAt)

	newCol, err := b.AddColumn("", in.Name, in.WipLimit, time.Now().UTC())
	if err != nil {
		return boardin.AddColumnOutput{}, translateDomainError(err)
	}

	colRecord := boardout.ColumnRecord{
		ID:         newCol.ID(),
		BoardID:    b.ID(),
		Name:       newCol.Name(),
		OrderIndex: newCol.OrderIndex(),
		WipLimit:   newCol.WipLimit(),
	}

	savedCol, err := s.boardRepo.AddColumn(ctx, b.ID(), colRecord)
	if err != nil {
		if errors.Is(err, boardout.ErrDuplicate) {
			return boardin.AddColumnOutput{}, fmt.Errorf("%w: %v", boardin.ErrDuplicateColumnName, err)
		}
		return boardin.AddColumnOutput{}, fmt.Errorf("add column: %w", err)
	}

	return boardin.AddColumnOutput{
		BoardID: b.ID(),
		Column: boardin.ColumnDTO{
			ID:         savedCol.ID,
			Name:       savedCol.Name,
			OrderIndex: savedCol.OrderIndex,
			WipLimit:   savedCol.WipLimit,
		},
	}, nil
}

