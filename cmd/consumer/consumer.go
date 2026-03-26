package consumer

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(kafkaCmd)
	Cmd.AddCommand(rabbitMQCmd)
	Cmd.AddCommand(natsCmd)
}

var Cmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start consumer",
}
