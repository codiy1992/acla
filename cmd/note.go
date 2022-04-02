package cmd

import (
	"aclt/anki"
	"aclt/dict"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type noteCmd struct {
	cmd        *cobra.Command
	deck       string
	model      string
	query      string
	fieldArray []string
	tags       []string
	override   bool
	remove     bool
	mainField  string
	dictQuery  bool
	fields     map[string]string
}

func (c noteCmd) newCmd() *cobra.Command {
	c.cmd = &cobra.Command{
		Use:   "note",
		Short: "Anki notes management",
		Long:  `Anki notes management`,
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			// split note ([]string) field:value into noteMap (map[string]string)
			c.fields = make(map[string]string)
			for _, str := range c.fieldArray {
				s := strings.Split(str, ":")
				if len(s) <= 1 {
					fmt.Println("Error: TODO ")
					c.cmd.Usage()
					os.Exit(1)
				}
				field := strings.Join(s[:1], "")
				value := strings.Join(s[1:], ":")
				c.fields[field] = value
			}

			err := c.checkFlags()
			if err != nil {
				fmt.Println("Error: " + err.Error())
				c.cmd.Usage()
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if c.remove {
				if len(c.tags) > 0 {
					// TODO removeTags
				} else {
					noteIds, _ := anki.Client().FindNotes(c.query)
					anki.Client().DeleteNotes(noteIds)
					return
				}
			} else {
				if c.query != "" && len(c.fieldArray) == 0 && len(c.tags) == 0 {
					noteIds, _ := anki.Client().FindNotes(c.query)
					notes, _ := anki.Client().NotesInfo(noteIds)
					json, err := json.Marshal(notes)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(string(json))
					return
				}

				if c.query == "" && len(c.fieldArray) > 0 {
					if c.withDict {
						if word, exists := c.fields["Word"]; exists {
							c.fields = dict.ToAnkiFields(word)
						}
					}
					// jsonData, _ := json.Marshal(c.fields)
					anki.AddNote(c.deck, c.model, c.fields, c.tags)
				}

				if c.query != "" {
					if len(c.fieldArray) > 0 || len(c.tags) > 0 {
						if c.withDict {
							if word, exists := c.fields["Word"]; exists {
								c.fields = dict.ToAnkiFields(word)
							}
						}
						anki.UpdateNote(c.query, c.fields, c.tags, c.override)
					}
				}
			}
		},
	}
	c.cmd.Flags().StringVarP(&c.deck, "deck", "d", "",
		"Deck name which notes located in")
	c.cmd.Flags().StringVarP(&c.model, "model", "m", "",
		"Model name (note type, required when executing addNote action)")
	c.cmd.Flags().StringVarP(&c.query, "query", "q", "",
		"Note query string like 'deck:English Word:present'(required when update or delete note)")
	c.cmd.Flags().StringArrayVarP(&c.fieldArray, "field", "f", []string{},
		"Note record (filed:value, can be specified multiple times)")
	c.cmd.Flags().StringArrayVarP(&c.tags, "tag", "t", []string{},
		"Note tags (can be specified multiple times)")
	c.cmd.Flags().BoolVar(&c.override, "override", false,
		"Override existing content, Valid only when update notes (default: false)")
	c.cmd.Flags().BoolVar(&c.remove, "remove", false,
		"Remove found notes or its tags (default: false, when specified --note will be ignored)")
	c.cmd.Flags().StringVarP(&c.mainField, "main-field", "M", "",
		"The main field that can be used to find specific note exactly "+
			"(Will be ignored when --query option is set)")
	c.cmd.Flags().BoolVar(&c.dictQuery, "dict-query", false,
		"")
	return c.cmd
}

func init() {
	c := noteCmd{}
	rootCmd.AddCommand(c.newCmd())
}

func (c noteCmd) checkFlags() error {
	if c.query == "" && len(c.fieldArray) == 0 {
		return errors.New("at least one of `--query` and `--field` must be specified")
	}
	if c.query == "" && len(c.fieldArray) > 0 {
		if c.mainField == "" {
			return errors.New("flag --main-field must be specified when --query is not set")
		}
		_, ok := c.fields[c.mainField]
		if !ok {
			return errors.New("main field must be set in note")
		}
	}

	if c.remove && c.query == "" {
		return errors.New("flag --query required when excuete removing action")
	}

	if !c.remove && c.query == "" && len(c.fieldArray) > 0 && c.model == "" {
		return errors.New("flag --model required when excueting addNote action")
	}

	if c.dictQuery && c.mainField != "Word" {
		return errors.New("withDict must work with xxxxxxx")
	}
	return nil
}
