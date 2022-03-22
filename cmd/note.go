package cmd

import (
	"aclt/anki"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var deck string
var model string
var query string
var fieldArray []string
var tags []string
var override bool
var remove bool
var mainField string
var fields map[string]string

func init() {
	rootCmd.AddCommand(noteCmd)
	noteCmd.Flags().StringVarP(&deck, "deck", "d", "",
		"Deck name which notes located in")
	noteCmd.Flags().StringVarP(&model, "model", "m", "",
		"Model name (note type, required when executing addNote action)")
	noteCmd.Flags().StringVarP(&query, "query", "q", "",
		"Note query string like 'deck:English Word:present'(required when update or delete note)")
	noteCmd.Flags().StringArrayVarP(&fieldArray, "field", "f", []string{},
		"Note record (filed:value, can be specified multiple times)")
	noteCmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{},
		"Note tags (can be specified multiple times)")
	noteCmd.Flags().BoolVar(&override, "override", false,
		"Override existing content, Valid only when update notes (default: false)")
	noteCmd.Flags().BoolVar(&remove, "remove", false,
		"Remove found notes or its tags (default: false, when specified --note will be ignored)")
	noteCmd.Flags().StringVarP(&mainField, "main-field", "M", "",
		"The main field that can be used to find specific note exactly "+
			"(Will be ignored when --query option is set)")

	noteCmd.PreRun = func(cmd *cobra.Command, args []string) {
		// split note ([]string) field:value into noteMap (map[string]string)
		fields = make(map[string]string)
		for _, str := range fieldArray {
			s := strings.Split(str, ":")
			if len(s) <= 1 {
				fmt.Println("Error: TODO ")
				noteCmd.Usage()
				os.Exit(1)
			}
			field := strings.Join(s[:1], "")
			value := strings.Join(s[1:], ":")
			fields[field] = value
		}

		err := checkFlags()
		if err != nil {
			fmt.Println("Error: " + err.Error())
			noteCmd.Usage()
			os.Exit(1)
		}
	}
}

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Anki notes management",
	Long:  `Anki notes management`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if remove {
			if len(tags) > 0 {
				// TODO removeTags
			} else {
				// TODO deleteNote
			}
		} else {
			if query != "" && len(fieldArray) == 0 && len(tags) == 0 {
				// TODO noteInfo
			}

			if query == "" && len(fieldArray) > 0 {
				anki.AddNote(deck, model, fields, tags)
			}

			if query != "" {
				if len(fieldArray) > 0 || len(tags) > 0 {
					anki.UpdateNote(query, fields, tags, override)
				}
			}
		}
	},
}

func checkFlags() error {
	if query == "" && len(fieldArray) == 0 {
		return errors.New("at least one of `--query` and `--field` must be specified")
	}
	if query == "" && len(fieldArray) > 0 {
		if mainField == "" {
			return errors.New("flag --main-field must be specified when --query is not set")
		}
		_, ok := fields[mainField]
		if !ok {
			return errors.New("main field must be set in note")
		}
	}

	if remove && query == "" {
		return errors.New("flag --query required when excuete removing action")
	}

	if !remove && query == "" && len(fieldArray) > 0 && model == "" {
		return errors.New("flag --model required when excueting addNote action")
	}
	return nil
}
