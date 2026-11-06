package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
)

type TaskCommands struct {
	createTask taskin.CreateTaskUseCase
	moveTask   taskin.MoveTaskUseCase
	listTasks  taskin.ListTasksUseCase
	deleteTask taskin.DeleteTaskUseCase
}

func NewTaskCommands(
	createTask taskin.CreateTaskUseCase,
	moveTask taskin.MoveTaskUseCase,
	listTasks taskin.ListTasksUseCase,
	deleteTask taskin.DeleteTaskUseCase,
) *TaskCommands {
	return &TaskCommands{
		createTask: createTask,
		moveTask:   moveTask,
		listTasks:  listTasks,
		deleteTask: deleteTask,
	}
}

func (c *TaskCommands) Execute(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required for task (create, move, list, delete)")
	}

	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("task create", flag.ContinueOnError)
		fs.SetOutput(out)
		boardID := fs.String("board-id", "", "Board ID (required)")
		column := fs.String("column", "Backlog", "Column ID or Name")
		title := fs.String("title", "", "Task title (required)")
		description := fs.String("description", "", "Task description")
		priority := fs.String("priority", "MEDIUM", "Task priority (LOW, MEDIUM, HIGH, URGENT)")
		dueDateStr := fs.String("due-date", "", "Due date (YYYY-MM-DD)")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *boardID == "" || *title == "" {
			return fmt.Errorf("--board-id and --title are required")
		}

		var dueDate *time.Time
		if *dueDateStr != "" {
			d, err := time.Parse("2006-01-02", *dueDateStr)
			if err != nil {
				return fmt.Errorf("invalid due-date format (use YYYY-MM-DD): %w", err)
			}
			dueDate = &d
		}

		res, err := c.createTask.Exec(ctx, taskin.CreateTaskInput{
			BoardID:     *boardID,
			Column:      *column,
			Title:       *title,
			Description: *description,
			Priority:    *priority,
			DueDate:     dueDate,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "Task created: ID %s [%s]\n", res.ID, res.Priority)
		return nil

	case "move":
		fs := flag.NewFlagSet("task move", flag.ContinueOnError)
		fs.SetOutput(out)
		taskID := fs.String("task-id", "", "Task ID (required)")
		toColumn := fs.String("to-column", "", "Destination column ID or Name (required)")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *taskID == "" || *toColumn == "" {
			return fmt.Errorf("--task-id and --to-column are required")
		}

		res, err := c.moveTask.Exec(ctx, taskin.MoveTaskInput{
			TaskID:           *taskID,
			TargetColumnID:   *toColumn,
			TargetColumnName: *toColumn,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "Task %s moved to '%s'.\n", res.TaskID, res.ColumnName)
		return nil

	case "list":
		fs := flag.NewFlagSet("task list", flag.ContinueOnError)
		fs.SetOutput(out)
		boardID := fs.String("board-id", "", "Board ID (required)")
		column := fs.String("column", "", "Optional column filter (ID or Name)")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *boardID == "" {
			return fmt.Errorf("--board-id is required")
		}

		res, err := c.listTasks.Exec(ctx, taskin.ListTasksInput{
			BoardID:    *boardID,
			ColumnID:   *column,
			ColumnName: *column,
		})
		if err != nil {
			return err
		}

		printTasksTable(out, res.Tasks)
		return nil

	case "delete":
		fs := flag.NewFlagSet("task delete", flag.ContinueOnError)
		fs.SetOutput(out)
		taskID := fs.String("task-id", "", "Task ID (required)")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *taskID == "" {
			return fmt.Errorf("--task-id is required")
		}

		_, err := c.deleteTask.Exec(ctx, taskin.DeleteTaskInput{
			TaskID: *taskID,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "Task %s deleted successfully.\n", *taskID)
		return nil

	default:
		return fmt.Errorf("unknown task subcommand: %s", args[0])
	}
}

