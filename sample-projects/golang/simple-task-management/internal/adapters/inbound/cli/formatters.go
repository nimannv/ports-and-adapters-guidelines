package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	boardin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/board"
	taskin "github.com/ports-and-adapters/simple-task-management/internal/core/application/ports/inbound/task"
)

func formatColumnsSummary(columns []boardin.ColumnDTO) string {
	parts := make([]string, len(columns))
	for i, c := range columns {
		if c.WipLimit != nil {
			parts[i] = fmt.Sprintf("%s (limit: %d)", c.Name, *c.WipLimit)
		} else {
			parts[i] = fmt.Sprintf("%s (unlimited)", c.Name)
		}
	}
	return strings.Join(parts, ", ")
}

func printTasksTable(w io.Writer, tasks []taskin.TaskItemDTO) {
	if len(tasks) == 0 {
		fmt.Fprintln(w, "No tasks found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(tw, "ID\tTITLE\tCOLUMN\tPRIORITY")
	for _, t := range tasks {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", t.ID, t.Title, t.ColumnName, t.Priority)
	}
	_ = tw.Flush()
}

