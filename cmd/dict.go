package cmd

import (
	"aclt/dict"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type dictCmd struct {
	cmd    *cobra.Command
	source string
	word   string
	format string
}

func (c dictCmd) newCmd() *cobra.Command {
	c.cmd = &cobra.Command{
		Use:   "dict",
		Short: "Dictionary query tool",
		Long:  `Dictionary query tool`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			result := dict.Query(c.word)
			if c.format != "yaml" {
				json, _ := json.Marshal(result)
				fmt.Println(string(json))
			} else {
				yaml, _ := yaml.Marshal(result)
				fmt.Println(string(yaml))
			}
		},
	}
	c.cmd.Flags().StringVarP(&c.word, "word", "w", "hello", "Specify word to query")
	c.cmd.Flags().StringVarP(&c.format, "format", "f", "yaml", "Output format (default: yaml)")
	return c.cmd
}

func init() {
	c := dictCmd{}
	rootCmd.AddCommand(c.newCmd())
}
