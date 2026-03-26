package server

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(restCmd)
}

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
}
