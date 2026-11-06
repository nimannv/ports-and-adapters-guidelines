package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
)

type BoardCommands struct {
	createBoard boardin.CreateBoardUseCase
	getBoard    boardin.GetBoardUseCase
	listBoards  boardin.ListBoardsUseCase
}

func NewBoardCommands(
	createBoard boardin.CreateBoardUseCase,
	getBoard boardin.GetBoardUseCase,
	listBoards boardin.ListBoardsUseCase,
) *BoardCommands {
	return &BoardCommands{
		createBoard: createBoard,
		getBoard:    getBoard,
		listBoards:  listBoards,
	}
}

func (c *BoardCommands) Execute(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required for board (create, get, list)")
	}

	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("board create", flag.ContinueOnError)
		fs.SetOutput(out)
		name := fs.String("name", "", "Board name (required)")
		description := fs.String("description", "", "Board description")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *name == "" {
			return fmt.Errorf("--name is required")
		}

		res, err := c.createBoard.Exec(ctx, boardin.CreateBoardInput{
			Name:        *name,
			Description: *description,
		})
		if err != nil {
			return err
		}

		fmt.Fprintln(out, "Board created:")
		fmt.Fprintf(out, "ID:          %s\n", res.ID)
		fmt.Fprintf(out, "Name:        %s\n", res.Name)
		fmt.Fprintf(out, "Columns:     %s\n", formatColumnsSummary(res.Columns))
		return nil

	case "get":
		fs := flag.NewFlagSet("board get", flag.ContinueOnError)
		fs.SetOutput(out)
		boardID := fs.String("board-id", "", "Board ID (required)")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *boardID == "" {
			return fmt.Errorf("--board-id is required")
		}

		res, err := c.getBoard.Exec(ctx, boardin.GetBoardInput{BoardID: *boardID})
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "Board ID:     %s\n", res.ID)
		fmt.Fprintf(out, "Name:         %s\n", res.Name)
		fmt.Fprintf(out, "Description:  %s\n", res.Description)
		fmt.Fprintf(out, "Columns:      %s\n", formatColumnsSummary(res.Columns))
		return nil

	case "list":
		res, err := c.listBoards.Exec(ctx, boardin.ListBoardsInput{})
		if err != nil {
			return err
		}
		if len(res.Boards) == 0 {
			fmt.Fprintln(out, "No boards found.")
			return nil
		}
		for _, b := range res.Boards {
			fmt.Fprintf(out, "ID: %s | %s (%s)\n", b.ID, b.Name, formatColumnsSummary(b.Columns))
		}
		return nil

	default:
		return fmt.Errorf("unknown board subcommand: %s", args[0])
	}
}
