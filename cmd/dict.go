package cmd

import (
	"aclt/dict"

	"github.com/spf13/cobra"
)

var source string

func init() {
	rootCmd.AddCommand(dictCmd)
	noteCmd.Flags().StringVarP(&source, "source", "s", "cambridge", "Specify directory source")
}

var dictCmd = &cobra.Command{
	Use:   "dict",
	Short: "Dictionary query tool",
	Long:  `Dictionary query tool`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// prevent scaffold nephew interactive
		dict.Query("interactive")
	},
}
