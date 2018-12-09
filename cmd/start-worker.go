package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startWorkerCmd represents the start-worker command
var startWorkerCmd = &cobra.Command{
	Use:   "start-worker",
	Short: "starts the worker",
	Long:  `starts the worker.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start-worker called")
	},
}

func init() {
	RootCmd.AddCommand(startWorkerCmd)
}
