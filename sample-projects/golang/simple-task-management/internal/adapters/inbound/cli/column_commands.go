package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
)

type ColumnCommands struct {
	addColumn boardin.AddColumnUseCase
}

func NewColumnCommands(addColumn boardin.AddColumnUseCase) *ColumnCommands {
	return &ColumnCommands{addColumn: addColumn}
}

func (c *ColumnCommands) Execute(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required for column (add)")
	}

	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("column add", flag.ContinueOnError)
		fs.SetOutput(out)
		boardID := fs.String("board-id", "", "Board ID (required)")
		name := fs.String("name", "", "Column name (required)")
		wipLimit := fs.Int("wip-limit", 0, "Optional WIP limit (>0)")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *boardID == "" || *name == "" {
			return fmt.Errorf("--board-id and --name are required")
		}

		var limitPtr *int
		if *wipLimit > 0 {
			limitPtr = wipLimit
		}

		res, err := c.addColumn.Exec(ctx, boardin.AddColumnInput{
			BoardID:  *boardID,
			Name:     *name,
			WipLimit: limitPtr,
		})
		if err != nil {
			return err
		}

		if res.Column.WipLimit != nil {
			fmt.Fprintf(out, "Column '%s' added with WIP limit %d.\n", res.Column.Name, *res.Column.WipLimit)
		} else {
			fmt.Fprintf(out, "Column '%s' added (unlimited).\n", res.Column.Name)
		}
		return nil

	default:
		return fmt.Errorf("unknown column subcommand: %s", args[0])
	}
}

