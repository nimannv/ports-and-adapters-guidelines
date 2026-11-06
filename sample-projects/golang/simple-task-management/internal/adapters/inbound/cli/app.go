package cli

import (
	"context"
	"fmt"
	"io"
)

type App struct {
	boardCmds  *BoardCommands
	columnCmds *ColumnCommands
	taskCmds   *TaskCommands
}

func NewApp(
	boardCmds *BoardCommands,
	columnCmds *ColumnCommands,
	taskCmds *TaskCommands,
) *App {
	return &App{
		boardCmds:  boardCmds,
		columnCmds: columnCmds,
		taskCmds:   taskCmds,
	}
}

func (a *App) Run(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		a.printHelp(out)
		return nil
	}

	command := args[0]
	subArgs := args[1:]

	switch command {
	case "board":
		return a.boardCmds.Execute(ctx, subArgs, out)
	case "column":
		return a.columnCmds.Execute(ctx, subArgs, out)
	case "task":
		return a.taskCmds.Execute(ctx, subArgs, out)
	default:
		return fmt.Errorf("unknown command %q. Use 'task-cli --help' for available commands", command)
	}
}

func (a *App) printHelp(out io.Writer) {
	helpText := `Simple Task Management CLI

Usage:
  task-cli <command> <subcommand> [flags]

Commands:
  board
    create   --name <string> [--description <string>]
    list
  column
    add      --board-id <uuid> --name <string> [--wip-limit <int>]
  task
    create   --board-id <uuid> [--column <name|uuid>] --title <string> [--description <string>] [--priority <LOW|MEDIUM|HIGH|URGENT>] [--due-date <YYYY-MM-DD>]
    move     --task-id <uuid> --to-column <name|uuid>
    list     --board-id <uuid> [--column <name|uuid>]
    delete   --task-id <uuid>
`
	fmt.Fprint(out, helpText)
}

