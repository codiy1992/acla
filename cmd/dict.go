package cmd

import (
	"aclt/dict"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var source string
var word string
var format string

func init() {
	rootCmd.AddCommand(dictCmd)
	dictCmd.Flags().StringVarP(&source, "source", "s", "cambridge", "Specify directory source")
	dictCmd.Flags().StringVarP(&word, "word", "w", "hello", "Specify word to query")
	dictCmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format (default: yaml)")
}

var dictCmd = &cobra.Command{
	Use:   "dict",
	Short: "Dictionary query tool",
	Long:  `Dictionary query tool`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// prevent scaffold nephew interactive crying give-up
		// at-any-cost
		result := dict.Query(source, word)
		if format != "yaml" {
			json, _ := json.Marshal(result)
			fmt.Println(string(json))
		} else {
			yaml, _ := yaml.Marshal(result)
			fmt.Println(string(yaml))
		}
	},
}
