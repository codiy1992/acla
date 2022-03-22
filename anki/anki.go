package anki

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/leonhfr/anki-connect-go"
)

var client anki.Client

func init() {
	client = *anki.NewDefaultClient()
}

func IsMainField(model string, fieldName string) bool {
	fields, err := client.ModelFieldNames(context.TODO(), model)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}
	if len(fields) == 0 {
		log.Fatalf("Model `%s` does not exist", model)
	}
	if fieldName != fields[0] {
		log.Printf("filed `%s` does match first field `%s` of model `%s`", fieldName, fields[0], model)
		return false
	}
	return true
}

func AddNote(deck string, model string, fields map[string]string, tags []string) int {

	for field, value := range fields {
		if strings.HasSuffix(field, "s") {
			fields[field] = "<ul><li>" + value + "</li></ul>"
		}
	}
	note := anki.NoteInput{
		Deck:   deck,
		Model:  model,
		Fields: fields,
		Options: map[string]interface{}{
			"allowDuplicate": false,
			"duplicateScope": "deck",
			"duplicateScopeOption": map[string]interface{}{
				"deckName":       deck,
				"checkChildren":  false,
				"checkAllModels": false,
			},
		},
		Tags: tags,
	}
	noteId, err := client.AddNote(context.TODO(), note)
	if err != nil {
		log.Fatalf("Create note failed: %s", err)
	}

	client.Suspend(context.TODO(), []int{noteId})
	return noteId
}

func UpdateNote(query string, fields map[string]string, tags []string, override bool) {

	notes, _ := client.FindNotes(context.TODO(), query)

	if len(notes) == 0 {
		log.Fatalf("No notes found!")
	}
	if len(notes) > 1 {
		log.Fatalf("More than one note found!")
	}

	if !override {
		res2, _ := client.NotesInfo(context.TODO(), notes)
		for field, value := range fields {
			if res2[0].Fields[field] == nil {
				log.Fatalf("field `%s` doesn't exists", field)
			}
			original := res2[0].Fields[field].(map[string]interface{})
			if strings.Trim(original["value"].(string), "\n ") == "" {
				fields[field] = "<ul><li>" + value + "</li></ul>"
			} else {
				re := regexp.MustCompile(`</li>[^<]*</ul>`)
				if re.MatchString(original["value"].(string)) {
					fields[field] = re.ReplaceAllString(
						original["value"].(string),
						"</li><li>"+value+"</li></ul>",
					)
				} else {
					fields[field] = strings.TrimRight(original["value"].(string), "\n") + "\n" + value
				}
			}
		}
	}

	note := anki.NoteFieldsInput{
		ID:     notes[0],
		Fields: fields,
	}

	if err := client.UpdateNote(context.TODO(), note); err != nil {
		log.Println("update failed:", err)
	}

	// TODO addTags
}
