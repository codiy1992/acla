package anki

import (
	"aclt/internal"
	"log"
	"regexp"
	"strings"

	"github.com/codiy1992/anki-connect-go"
)

var client anki.Client

func init() {
	client = *anki.NewDefaultClient()
}

func Client() *anki.Client {
	return &client
}

func AddVocab(deck string, model string, fields map[string]string, tags []string) int {

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
	noteId, err := client.AddNote(note)
	if err != nil {
		log.Fatalf("Create note failed: %s", err)
	}
	//
	internal.Put(fields)
	return noteId
}

func UpdateVocab(query string, fields map[string]string, tags []string, override bool) {

	notes, _ := client.FindNotes(query)
	var primaryFieldVaule string
	if len(notes) == 0 {
		log.Fatalf("No notes found!")
	}
	if len(notes) > 1 {
		log.Fatalf("More than one note found!")
	}

	if !override {
		res, _ := client.NotesInfo(notes)
		if primaryField, exists := res[0].Fields["Word"]; exists {
			primaryFieldVaule = primaryField.(map[string]interface{})["value"].(string)
		}
		for field, value := range fields {
			if res[0].Fields[field] == "" {
				log.Fatalf("field `%s` doesn't exists", field)
			}
			original := res[0].Fields[field].(map[string]interface{})
			fields[field] = strings.TrimRight(original["value"].(string), "\n") + "\n" + value
			if strings.Trim(original["value"].(string), "\n ") == "" {
				re := regexp.MustCompile(`<\s*img[^>]*>`)
				if !re.MatchString(value) {
					fields[field] = "<ul><li>" + value + "</li></ul>"
				}
			} else {
				re := regexp.MustCompile(`</li>[^<]*</ul>`)
				if re.MatchString(original["value"].(string)) {
					fields[field] = re.ReplaceAllString(
						original["value"].(string),
						"</li>\n<li>"+value+"</li>\n</ul>",
					)
				}
			}
		}
	}

	note := anki.NoteFieldsInput{
		ID:     notes[0],
		Fields: fields,
	}

	if err := client.UpdateNote(note); err != nil {
		log.Println("update failed:", err)
	}
	fields["Word"] = primaryFieldVaule
	internal.Put(fields)
	// TODO addTags
}
