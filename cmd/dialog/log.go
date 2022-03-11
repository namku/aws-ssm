package dialog

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/spf13/cobra"
)

// severity: ["Error", "Warning", "Info"]
// msg: message to print
func Log(severity string, msg string, cmd *cobra.Command) {
	switch severity {
	case "Error":
		colorstring.Println("\n[bold][red]Error: [reset]" + msg)
		fmt.Println()
		cmd.Help()
		fmt.Println()
		os.Exit(1)
	case "Warning":
		colorstring.Println("\n[bold][yellow]Warning: [reset]" + msg)
		fmt.Println()
	case "Info":
		colorstring.Println("\n[bold][green]" + msg + "[reset]")
	}
}
