package internal

import (
	"context"
	"log"

	"github.com/leonhfr/anki-connect-go"
)

var client anki.Client

func IsFirstField(model string, fieldName string) bool {
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

func CreateNote(deck string, model string, fieldName string, fieldValue string, tags []string) int {
	note := anki.NoteInput{
		Deck:  deck,
		Model: model,
		Fields: map[string]interface{}{
			fieldName: fieldValue,
		},
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

func UpdateFiled(query string, fieldName string, fieldValue string, append bool) {

	notes, _ := client.FindNotes(context.TODO(), query)

	if len(notes) == 0 {
		log.Fatalf("No notes found!")
	}
	if len(notes) > 1 {
		log.Fatalf("More than one note found!")
	}

	fields := map[string]interface{}{
		fieldName: fieldValue,
	}

	if append {
		res2, _ := client.NotesInfo(context.TODO(), notes)

		if res2[0].Fields[fieldName] == nil {
			log.Fatalf("field `%s` doesn't exists", fieldName)
		}
		original := res2[0].Fields[fieldName].(map[string]interface{})

		fields[fieldName] = original["value"].(string) + fieldValue
	}

	note := anki.NoteFieldsInput{
		ID:     notes[0],
		Fields: fields,
	}

	if err := client.UpdateNote(context.TODO(), note); err != nil {
		log.Println("update failed:", err)
	}
}

func init() {
	client = *anki.NewDefaultClient()
}
